/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.vmware.client.automation.vcuilib.commoncode.GlobalFunction;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.UiComponentProperty;
import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.componentframework.controls.mx.Label;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * This class handles the work with Property Grid (usually on Ready
 * to Complete page of wizards) component
 */
public class PropertyGridControl {

    /**
     * The type of the selected item.
     */
    public static enum ItemType {
        TEXT, LINK
    }

    /**
     * This class represents a property item.
     */
    public static class PropertyGridItem {
        private ItemType type;
        private String uid;

        public PropertyGridItem(final String uid, final ItemType type) {
            this.uid = uid;
            this.type = type;
        }

        public ItemType getType() {
            return type;
        }

        public void setType(ItemType type) {
            this.type = type;
        }

        public String getUid() {
            return uid;
        }

        public void setUid(String uid) {
            this.uid = uid;
        }
    }

    /**
     * Interface for parsers which should provide logic to parse individual items.
     */
    public static interface PropertyGridItemParser {
        public PropertyGridItem parseItem(final String itemId);
    }

    /**
     * Represents a row in the property grid.
     */
    public static final class PropertyGridRow {
        private final String rowId;
        private Label label = null;
        private final List<String> propertyGridItemIdList = new ArrayList<String>();

        public PropertyGridRow(final String gridId) {
            this.rowId = gridId;
            parseProperties();
        }

        private void parseProperties() {
            final List<List<String>> allProperties = getChildProperties(rowId);

            for (int i = 0; i < allProperties.get(UID_PROPERTIES_INDEX).size(); i++) {
                if (Boolean.parseBoolean(allProperties.get(VISIBLE_PROPERTIES_INDEX).get(i))) {
                    final String className = allProperties.get(CLASS_NAME_PROPERTIES_INDEX).get(i);
                    final String uid = new StringBuilder().append("uid=")
                        .append(allProperties.get(UID_PROPERTIES_INDEX).get(i)).toString();

                    if (className.equals(LABEL_GRID_ITEM_CLASS)) {
                        label = new Label(uid, BrowserUtil.flashSelenium);
                    } else if (className.equals(PROPERTY_GRID_ITEM_CLASS)) {
                        propertyGridItemIdList.add(uid);
                    }
                }
            }
        }

        public Label getLabel() {
            return label;
        }

        public List<String> getPropertyGridItemIdList() {
            return propertyGridItemIdList;
        }
    }

    private static final String PROPERTY_DELIMETER = "^^";
    private static final String PROPERTY_SPLITTER = "\\^\\^";
    private static final String MULTI_PROPERTY_DEFAULT_DELIMETER = ",";
    private static final String ALL_PROPERTIES = "[].uid,[].className,[].visible";
    private static final String PROPERTY_GRID_ROW_CLASS = "PropertyGridRow";
    private static final String LABEL_GRID_ITEM_CLASS = "LabelGridItem";
    private static final String PROPERTY_GRID_ITEM_CLASS = "PropertyGridItem";

    private static final int UID_PROPERTIES_INDEX = 0;
    private static final int CLASS_NAME_PROPERTIES_INDEX = 1;
    private static final int VISIBLE_PROPERTIES_INDEX = 2;
    private static final int WAIT_FOR_VISIBLE_ITEM_RETRIES = 30;

    // The id of the grid.
    private final String gridId;

    // A list containing this grid's rows.
    protected final List<PropertyGridRow> gridRows = new ArrayList<PropertyGridRow>();

    /**
     * Creating a PropertyGrid object implies that the property grid is visible on screen.
     *
     * @param gridId
     */
    public PropertyGridControl(final String gridId) {
        this.gridId = gridId;
        GlobalFunction.waitforUIComponentVisible(this.gridId, WAIT_FOR_VISIBLE_ITEM_RETRIES);
        initPropertyGrid();
    }

    /**
     * Initialize the property grid.
     */
    private void initPropertyGrid() {
        final List<List<String>> allProperties = getChildProperties(gridId);

        for (int i = 0; i < allProperties.get(UID_PROPERTIES_INDEX).size(); i++) {
            if (allProperties.get(CLASS_NAME_PROPERTIES_INDEX).get(i).equals(PROPERTY_GRID_ROW_CLASS)
                && Boolean.parseBoolean(allProperties.get(VISIBLE_PROPERTIES_INDEX).get(i))) {
                final StringBuilder itemId = new StringBuilder().append(gridId).append("/uid=")
                    .append(allProperties.get(UID_PROPERTIES_INDEX).get(i));
                gridRows.add(new PropertyGridRow(itemId.toString()));
            }
        }
    }

    /**
     * Get the visible child properties of the given component.
     *
     * @param componentId
     * @return
     */
    private static List<List<String>> getChildProperties(final String componentId) {
        final List<List<String>> items = new ArrayList<List<String>>();

        final DisplayObject object = new DisplayObject(componentId, BrowserUtil.flashSelenium);
        final String rawValues = object.getProperties(ALL_PROPERTIES, PROPERTY_DELIMETER);
        final String[] allValues = rawValues.split(MULTI_PROPERTY_DEFAULT_DELIMETER);

        items.add(Arrays.asList(allValues[UID_PROPERTIES_INDEX].split(PROPERTY_SPLITTER)));
        items.add(Arrays.asList(allValues[CLASS_NAME_PROPERTIES_INDEX].split(PROPERTY_SPLITTER)));
        items.add(Arrays.asList(allValues[VISIBLE_PROPERTIES_INDEX].split(PROPERTY_SPLITTER)));

        return items;
    }

    /**
     * Number of rows for this grid.
     *
     * @return
     */
    public int getNumberOfRows() {
        return gridRows.size();
    }

    /**
     * Get a list of all of the rows.
     *
     * @return
     */
    public List<PropertyGridRow> getAllRows() {
        return gridRows;
    }

    /**
     * Parse property grid items with the given parser.
     *
     * @param parser
     * @return
     */
    public Map<String, List<PropertyGridItem>> parseRowItems(final PropertyGridItemParser parser) {
        final Map<String, List<PropertyGridItem>> propertyMapping = new HashMap<String, List<PropertyGridItem>>();

        for (PropertyGridRow row : gridRows) {
            final List<String> itemList = row.getPropertyGridItemIdList();
            final List<PropertyGridItem> parsedItems = new ArrayList<PropertyGridControl.PropertyGridItem>();

            for (String item : itemList) {
                parsedItems.add(parser.parseItem(item));
            }

            propertyMapping.put(row.getLabel().getProperty(UiComponentProperty.LABEL.getName()).trim(), parsedItems);
        }

        return propertyMapping;
    }
}
