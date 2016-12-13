/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Set;

import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;

import org.apache.commons.lang.StringUtils;

/**
 * Utility for finding components and their properties - use only in control wrappers.
 */
public class PropertiesUtil {
    private static final String NAME_PROPERTY = "name";
    private static final String ARRAY_DELIMITER = ";;;";
    private static final String PROPERTIES_DELIMITER = "!!!";
    private static final String PROPERTIES_QUERY_DELIMITER = ",";

    /**
     * Finds the first-level children of element with specified id that suit the provided filters
     *
     * @param String componentId - id of the component to be searched
     * @param Map<property name, property value> filters - filters to be applied when looking for suitable children
     * @return List<String> - list of the names of the suitable children
     */
    public static List<String> findChildrenNames(String componentId, Map<String, String> filters) {
        return findChildrenNames(componentId, filters, 1);
    }

    /**
     * Finds the children of element with specified id that are on the specified depth and suit
     * the provided filters
     *
     * @param String componentId - id of the component to be searched
     * @param Map<String, String> filters - filters to be applied when looking for suitable children
     * @param int level - desired depth of the children
     * @return List<String> - list of the names of the suitable children
     */
    public static List<String> findChildrenNames(String componentId, Map<String, String> filters, int level) {
        Map<String, String> childrenPropertiesFilters = new HashMap<String, String>();
        for (Entry<String, String> filter : filters.entrySet()) {
            childrenPropertiesFilters.put(getChildPropertyId(filter.getKey(), level), filter.getValue());
        }
        return findNames(componentId, childrenPropertiesFilters, level);
    }

    /**
     * Get properties for component with specific ID
     *
     * @param componentId
     * @param properties
     * @return Map<PropertyName, values> - if the the property has single value, the array will have one element
     */
    public static Map<String, String[]> getProperties(String componentId, String... properties) {
        UIComponent component = new UIComponent(componentId, BrowserUtil.flashSelenium);
        String query = StringUtils.join(properties, PROPERTIES_QUERY_DELIMITER);

        String rawResult = component.getProperties(query, ARRAY_DELIMITER, PROPERTIES_DELIMITER);
        String[] propertiesValues = rawResult.split(PROPERTIES_DELIMITER, -1);
        if (propertiesValues.length != properties.length) {
            return null;
        }

        Map<String, String[]> result = new HashMap<String, String[]>();
        for (int i = 0; i < properties.length; i++) {
            result.put(properties[i], propertiesValues[i].split(ARRAY_DELIMITER, -1));
        }
        return result;
    }

    private static List<String> findNames(String componentId, Map<String, String> filters, int level) {
        String nameProperty = getChildPropertyId(NAME_PROPERTY, level);

        List<String> result = new ArrayList<String>();

        Set<String> inspectedProperties = new HashSet<>(filters.keySet());
        inspectedProperties.add(nameProperty);

        Map<String, String[]> propertyValues = getProperties(
            componentId,
            inspectedProperties.toArray(new String[] {}));

        if (propertyValues == null || !propertyValues.containsKey(nameProperty)) {
            return result;
        }

        int resultsCount = propertyValues.get(nameProperty).length;
        for (int i = 0; i < resultsCount; i++) {
            Map<String, String> componentProperties = extractPropertiesForComponent(propertyValues, i);
            String componentName = componentProperties.get(nameProperty);
            componentProperties.remove(nameProperty);
            if (propertiesMatch(componentProperties, filters)) {
                result.add(componentName);
            }
        }

        return result;
    }

    private static boolean propertiesMatch(Map<String, String> left, Map<String, String> right) {
        for (Entry<String, String> entry : left.entrySet()) {
            if (!entry.getValue().equals(right.get(entry.getKey()))) {
                return false;
            }
        }

        return true;
    }

    private static Map<String, String> extractPropertiesForComponent(Map<String, String[]> propertyValues,
        int componentIndex) {

        Map<String, String> properties = new HashMap<>();
        for (Entry<String, String[]> entry : propertyValues.entrySet()) {
            properties.put(entry.getKey(), entry.getValue()[componentIndex]);
        }
        return properties;
    }

    private static String getChildPropertyId(String property, int level) {
        StringBuilder childPropertyId = new StringBuilder();
        for (int i = 0; i < level; i++) {
            childPropertyId.append("[].");
        }
        childPropertyId.append(property);
        return childPropertyId.toString();
    }
}

