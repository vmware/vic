/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.control;

import com.vmware.flexui.componentframework.DisplayObjectContainer;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.client.automation.components.control.PropertyGridControl;
import com.vmware.client.automation.components.control.PropertyGridControl.PropertyGridItem;
import com.vmware.client.automation.components.control.PropertyGridControl.PropertyGridItemParser;


/**
 * Parser for Property Grid items when only one text field is expected in the right container of the grid row
 */
public final class UiTextFieldPropertyGridItemParser implements PropertyGridItemParser {
    private static final String NAME_PROPERTY = "name";
    private static final String UITEXTFIELD_CLASS_NAME = "/className=UITextField";
    private static UiTextFieldPropertyGridItemParser parserInstance = null;

    private UiTextFieldPropertyGridItemParser() {
        // Make the constructor private to disallow external object creation.
    }

    public static UiTextFieldPropertyGridItemParser getInstance() {
        if (parserInstance == null) {
            parserInstance = new UiTextFieldPropertyGridItemParser();
        }
        return parserInstance;
    }

    @Override
    public PropertyGridItem parseItem(String itemId) {
        PropertyGridItem item = null;

        DisplayObjectContainer object = new DisplayObjectContainer(itemId, BrowserUtil.flashSelenium);
        String uiTextFieldId = NAME_PROPERTY + "=" + object.getProperty(NAME_PROPERTY) + UITEXTFIELD_CLASS_NAME;
        item = new PropertyGridItem(uiTextFieldId, PropertyGridControl.ItemType.TEXT);

        return item;
    }

}
