package com.vmware.suitaf.apl;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import com.vmware.suitaf.apl.IDPart.MatchType;

/**
 * This class is an element of the general component ID format. An ID group
 * holds the information needed for the identification of one screen component.
 * It is in the form of Property-Value pairs that are used to match against the
 * properties of the screen components. Besides the {@link Category#REGULAR}
 * Property-Value pairs, there are some {@link Category#SYSTEM} pairs. The
 * system pairs are not matched directly to component's properties but
 * control the usage of the ID group.
 * <br><br>
 * The most important of the system properties is the {@link
 * Property#GROUPROLE}. Each ID group has one such Property-Value pair that
 * determines if this is the main ID group or it has other role. The
 * available ID group roles are enlisted in the constant container class:
 * {@link IDPart.GROUPROLE}
 * <br><br>
 * <b>NOTE:</b> This class is designed to be immutable. Please do not change
 * this feature.
 * <br><br>
 * @author dkozhuharov
 */
public final class IDGroup {

    /** The map of parts in this ID group. */
    private final Map<Property, List<String>> parts;

    private IDGroup(IDGroup old, IDPart[] partArray) {
        parts = new HashMap<Property, List<String>>();

        if (old != null) {
            // If an old instance is given copy its content
            for (Property property : old.parts.keySet()) {
                parts.put(property,
                        new ArrayList<String>(old.parts.get(property)));
            }
        }
        // Populate all passed IDParts in the parts map
        for (IDPart p : partArray) {
            putIDPart(p);
        }
        // Set the default group role
        if (!has(Property.GROUPROLE)) {
            putIDPart(IDPart.GROUPROLE.MAIN);
        }
        // Set default visibility filter
        if (!has(Property.VISIBLE)) {
            putIDPart(IDPart.from(
                    Property.VISIBLE, Boolean.TRUE.toString()));
        }
        else {
            // Remove visibility filter if both true/false are allowed
            if (parts.get(Property.VISIBLE)
                    .contains(Boolean.TRUE.toString())
                    && parts.get(Property.VISIBLE)
                    .contains(Boolean.FALSE.toString())) {
                parts.remove(Property.VISIBLE);
            }
        }
    }

    private void putIDPart(IDPart part) {
        // Skip null id-parts
        if (part == null) {
            return;
        }
        if (!parts.containsKey(part.property)) {
            parts.put(part.property, new ArrayList<String>());
        }
        parts.get(part.property).add(part.value);
    }

    /**
     * Method checks if this ID group has a Property-Value pair for a given
     * {@link Property}
     * @param checkProperty - property to be checked
     * @return <b>true</b> if a corresponding Property-Value pair is found
     */
    public boolean has(Property checkProperty) {
        return parts.containsKey(checkProperty);
    }

    /**
     * Method checks if this ID group has a particular Property-Value pair.
     * The pair is given as {@link IDPart}.
     * @param checkPart - the Property-Value pair to be checked
     * @return <b>true</b> if the Property-Value pair is found
     */
    public boolean has(IDPart checkPart) {
        return parts.containsKey(checkPart.property) &&
                getValues(checkPart.property).contains(checkPart.value);
    }

    /**
     * Accessor method that retrieves a set of all properties for which there
     * is data in this ID group.
     *
     * @return a set of {@link Property} instances
     */
    public Set<Property> getProperties() {
        return parts.keySet();
    }

    /**
     * Accessor method that retrieves the first value stored for a given
     * {@link Property}. For advanced property filters returns only the first
     * parameter as value.
     * @param property - the property for which the value will be retrieved
     * @return - the string value if found or <b>null</b> if no data for
     * this property was stored.
     */
    public String getValue(Property property) {
        List<String> values = parts.get(property);
        if (values != null && values.size() > 0) {
            String temp = values.get(0);
            String[] filterValues = IDPart.MatchType.unpackArgs(temp);
            return filterValues[0];
        }
        return null;
    }

    /**
     * Accessor method that retrieves the list of values stored for a given
     * {@link Property}.
     * @param property - the property for which the values will be retrieved
     * @return - the list of string values if found or <b>null</b> if no data
     * for this property was stored.
     */
    public List<String> getValues(Property property) {
        return parts.get(property);
    }

    // ===================================================================
    // Public static factory for identifiers
    // ===================================================================
    private Number toNumber(String value) {
        try {
            return Long.valueOf(value);
        }
        catch (Exception e) {
            try {
                return Double.valueOf(value);
            }
            catch (Exception e1) {
                return null;
            }
        }

    }
    /**
     * Match a single-valued component property against the id-part in this
     * id-group corresponding to that property.
     * @param componentValue  - the property value of the component
     * @param propertyToMatch - the {@link Property} identifier from the
     * single-valued properties to be matched
     * @return <b>true</b> if the match criterion has been met
     */
    public boolean matchComponentValue(
            String componentValue, Property propertyToMatch) {
        List<String> matchValues = parts.get(propertyToMatch);
        if (matchValues == null) {
            return false;
        }
        // This cycle iterates amongst the alternative property-matching values
        // and tries to match the property-value with any of them
        for (String matchValue : matchValues) {
            if (componentValue == null) {
                if (matchValue == null) {
                    // Support matching with "null"
                    return true;
                }
                continue;
            }

            MatchType mt = MatchType.unpackType(matchValue);
            String[] args = MatchType.unpackArgs(matchValue);

            if (mt.isNumeric) {
                Number pr = toNumber(componentValue);
                Number cmp0 = args.length>0? toNumber(args[0]): null;
                Number cmp1 = args.length>1? toNumber(args[1]): null;

                // Check if no necessary numbers are passed as parameters
                // to the numeric property filter
                if (mt == MatchType.N_EQUALS) {
                    if (cmp0 == null || cmp1 == null) {
                        continue;
                    }
                }
                else if (mt == MatchType.N_FITS_INCL || mt == MatchType.N_FITS_EXCL) {
                    if (cmp0 == null && cmp1 == null) {
                        continue;
                    }
                }

                // Numeric property filter in Long format
                if (cmp0 instanceof Long || cmp1 instanceof Long) {
                    long prL = pr.longValue();
                    long cmp0L = (cmp0 != null? cmp0.longValue(): 0);
                    long cmp1L = (cmp1 != null? cmp1.longValue(): 0);
                    if (mt == MatchType.N_EQUALS &&
                            (cmp0L - cmp1L <= prL) &&
                            (prL <= cmp0L + cmp1L) ) {
                        return true;
                    }
                    else if (mt == MatchType.N_FITS_INCL &&
                            (cmp0 == null || cmp0L <= prL) &&
                            (cmp1 == null || prL <= cmp1L) ) {
                        return true;
                    }
                    else if (mt == MatchType.N_FITS_EXCL &&
                            (cmp0 == null || cmp0L < prL) &&
                            (cmp1 == null || prL < cmp1L) ) {
                        return true;
                    }
                }
                // Numeric property filter in Double format
                else {
                    double prD = pr.doubleValue();
                    double cmp0D = (cmp0 != null? cmp0.doubleValue(): 0);
                    double cmp1D = (cmp1 != null? cmp1.doubleValue(): 0);
                    if (mt == MatchType.N_EQUALS &&
                            (cmp0D - cmp1D <= prD) &&
                            (prD <= cmp0D + cmp1D) ) {
                        return true;
                    }
                    else if (mt == MatchType.N_FITS_INCL &&
                            (cmp0 == null || cmp0D <= prD) &&
                            (cmp1 == null || prD <= cmp1D) ) {
                        return true;
                    }
                    else if (mt == MatchType.N_FITS_EXCL &&
                            (cmp0 == null || cmp0D < prD) &&
                            (cmp1 == null || prD < cmp1D) ) {
                        return true;
                    }
                }

            }
            else {
                // String property filters
                if (mt == MatchType.S_EQUALS &&
                        componentValue.equals(args[0])) {
                    return true;
                }
                else if (mt == MatchType.S_STARTS &&
                        componentValue.startsWith(args[0])) {
                    return true;
                }
                else if (mt == MatchType.S_ENDS &&
                        componentValue.endsWith(args[0])) {
                    return true;
                }
                else if (mt == MatchType.S_CONTAINS &&
                        componentValue.contains(args[0])) {
                    return true;
                }
                else if (mt == MatchType.S_MATCHES &&
                        componentValue.matches(args[0])) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * Match a list-valued component property against the id-part in this
     * id-group corresponding to that property.
     * @param componentValues - the list of property value of the component
     * @param propertyToMatch - the property identifier of the list-valued
     * property to be matched
     * @return <b>true</b> if the match criterion has been met
     */
    public boolean matchComponentValueList(
            List<String> componentValues, Property propertyToMatch) {
        if (componentValues != null) {
            for (String componentValue : componentValues) {
                if (matchComponentValue(componentValue, propertyToMatch))
                    return true;
            }
        }
        return false;
    }

    /**
     * Match a grid-valued component property against the id-part in this
     * id-group corresponding to that property.
     * @param componentValues - the grid of property value of the component
     * @param propertyToMatch - the property identifier of the grid-valued
     * property to be matched
     * @return <b>true</b> if the match criterion has been met
     */
    public boolean matchComponentValueGrid(
            List<List<String>> componentValues, Property propertyToMatch) {
        if (componentValues != null) {
            for (List<String> componentValueList : componentValues) {
                if (matchComponentValueList(componentValueList, propertyToMatch))
                    return true;
            }
        }
        return false;
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("IDGroup").append(parts.get(Property.GROUPROLE)).append(" {");
        for (Property prop : parts.keySet()) {
            if (!prop.equals(Property.GROUPROLE)) {
                sb.append(prop).append("=").append(parts.get(prop)).append(", ");
            }
        }
        sb.append("}");
        return sb.toString();
    }

    // ===================================================================
    // Public static factory for identifiers
    // ===================================================================
    private static final String GRAPCHICAL_SUFFIXES = "(.*bmp$)|(.*jpg$)|(.*jpeg$)|(.*png$)";
    /**
     * This method serves as an universal translator of identifiers to the
     * {@link IDGroup} format used in the APL interface. Accepted are the
     * following inputs:
     * <li> <b>null</b> - returns null - transparent for nulls
     * <li> {@link IDGroup} - returns the same value - transparent for already
     * converted identifiers. If accompanied by one or more property filters in
     * the form of {@link IDPart}s, then a new {@link IDGroup} instance is
     * created. It contains the data of the original IDGroup plus added the new
     * property filters.
     * <li> {@link String} - allows creation of ad-hoc identifiers at the place
     * of function call. If the string starts with "//" or "./" prefix - it is
     * treated as XPath, otherwise it is accepted for DirectID. The string ID
     * could be accompanied by one or more property filters in the form of
     * {@link IDPart}s. They are added to the created {@link IDGroup} instance.
     * <li> {@link Object} - any object can be a link to a IDGroup identifier
     * if registered with the {@link IDConstants#
     * setID(Object, Object, IDPart...)} static method
     * <br><br>
     * <b>NOTE:</b> If no {@link IDPart} element of type {@link
     * Property#GROUPROLE} is give - the {@link IDPart.GROUPROLE#MAIN} is put
     * by default.
     * <br><br>
     * <b>NOTE:</b> If no {@link IDPart} element of type {@link
     * Property#VISIBLE} is give - the {@link Property#VISIBLE}->true is put
     * by default.
     * <br><br>
     * @param id - identifier to be converted
     * @param propertyFilters - optional property filters to be added
     * @return - the resulting {@link IDGroup} format identifier
     * @throws RuntimeException if unsupported format was provided for ID
     */
    public static IDGroup toIDGroup(Object id, IDPart...propertyFilters) {
        if (id == null) {
            return null;
        }
        else if (id instanceof IDGroup) {
            if (propertyFilters.length > 0) {
                return new IDGroup((IDGroup) id, propertyFilters);
            }
            else {
                return (IDGroup) id;
            }
        }
        else if (id instanceof String) {
            String stringID = (String) id;
            List<IDPart> parts = new ArrayList<IDPart>();
            // Add the primary identifier
            if(stringID.matches(GRAPCHICAL_SUFFIXES)){
                parts.add(IDPart.from(Property.GRAPHICAL_ID, stringID));
            }else{
                parts.add(IDPart.id(stringID));
            }
            // Adding ad-hoc property filters if provided
            parts.addAll(Arrays.asList(propertyFilters));
            return new IDGroup(null, parts.toArray(new IDPart[0]));
        }
        throw new RuntimeException("Unsupported conversion type: " +
                id.getClass().getCanonicalName());
    }

}
