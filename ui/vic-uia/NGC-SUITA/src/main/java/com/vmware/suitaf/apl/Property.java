package com.vmware.suitaf.apl;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.regex.Pattern;

import com.vmware.hsua.common.datamodel.AbstractProperty;
import com.vmware.suitaf.apl.IDPart.MatchType;

/**
 * This enum is a general public list of the properties which are mostly used.
 * Each implementor should do its own mapping between this enum set and the
 * specific property names in their implementation domain.
 * <br><br>
 * Also these properties are attribute of the identifier parts
 * (see {@link IDPart}). Identifier parts are the base of the general ID format
 * and are used in the construction of screen component identifiers.
 * <br><br>
 * @author dkozhuharov
 */
public enum Property {

    // =====================================================================
    // Regular properties
    // =====================================================================

    /** The caption of the component */
    CAPTION ("The caption of the component - gets label property of the component",
            Category.REGULAR, Category.INFO),

    /** The text that stays behind a component.
     * Usually corresponds to the text, label or title property of a component */
    TEXT ("The text that stays behind a component",
            Category.REGULAR, Category.BRIEF, Category.INFO),

    /** The text shows component input errors */
    ERRORTEXT ("The text shows component input errors",
            Category.REGULAR, Category.INFO),

    /** The native name of the component's class */
    CLASS_NAME ("The native name of the component's class",
            Category.REGULAR),

    /** Boolean property showing if the component is visible */
    VISIBLE ("Boolean property showing if the component is visible",
            Category.REGULAR, Category.ACCESSIBILITY),

    /** Boolean property showing if the component is enabled for input */
    ENABLED ("Boolean property showing if the component is enabled for input",
            Category.REGULAR, Category.ACCESSIBILITY),

    /** Boolean property showing if the component has the key-type focus */
    FOCUSED ("Boolean property showing if the component has the key-type focus",
            Category.REGULAR, Category.SELECTION),

    /** Integer property showing the horizontal offset of the component */
    POS_X ("Integer property showing the horizontal offset of the component",
            Category.REGULAR),

    /** Integer property showing the vertical offset of the component */
    POS_Y ("Integer property showing the vertical offset of the component",
            Category.REGULAR),

    /** Integer property showing the vertical size of the component */
    POS_HEIGHT ("Integer property showing the vertical size of the component",
            Category.REGULAR, Category.POSITIONING),

    /** Integer property showing the horizontal size of the component */
    POS_WIDTH ("Integer property showing the horizontal size of the component",
            Category.REGULAR, Category.POSITIONING),

    /** Current value of a value container.
     * Implementation must support retrieval for the following generic types:
     * <li> text-box         - the text value
     * <li> numeric stepper  - the numeric value
     * <li> date-time picker - the date-time value
     * <li> combo-box - the selected combo-box item
     * <li> check-box - the true(checked)/false(unchecked) state value
     * <li> list      - the selected list item
     * <li> radio-button - the true(selected)/false(unselected) state value
     * <li> tab-control  - the selected tab-name value
     * <li> tree-control - the selected tree-item value
     *
     */
    VALUE ("Current value of a value container",
            Category.REGULAR, Category.INFO),

    /** Index of the current value of a value container.
     * Implementation must support retrieval for the following generic types:
     * <li> combo-box
     * <li> grid
     * <li> list
     * <li> tab-control
     * <li> tree-control
     */
    VALUE_INDEX ("Index of the current value of a value container",
            Category.REGULAR, Category.INFO),

    /** List of available values of a value container.
     * Implementation must support retrieval for the following generic types:
     * <li> combo-box
     * <li> grid
     * <li> list
     * <li> tab-control
     * <li> tree-control
     */
    VALUE_LIST ("List of available values of a value container",
            Category.REGULAR, Category.INFO),

    /** Grid column names */
    GRID_COLUMNS ("Grid column names",
            Category.REGULAR, Category.TABULAR),

    /** Grid data content */
    GRID_DATA ("Grid data content",
            Category.REGULAR, Category.TABULAR),

    /** Child component identifiers */
    CHILD_IDS ("Child component identifiers",
            Category.REGULAR),

    /** Automation Name identifier */
    AUTOMATION_NAME ("Automation Name identifier",
            Category.REGULAR, Category.BRIEF, Category.ID),

    /** ID identifier */

            ID ("ID identifier",
            Category.REGULAR, Category.BRIEF, Category.ID),

    /** Number of children */
    CHILDREN_NUMBER ("Number of children",
            Category.REGULAR, Category.INFO),

    /** Current state */
    CURRENT_STATE ("Current state",
            Category.REGULAR, Category.INFO),

    /** The show progress status */
    SHOW_PROGRESS("The showProgress status of a data-grid",
            Category.REGULAR, Category.INFO),

    /** The show progress status */
    TOOLTIP("The toolTip of a UI element",
            Category.REGULAR, Category.INFO),

    // =====================================================================
    // System properties
    // =====================================================================

    /** XPath format identifier. */
    X_PATH_ID ("XPath format identifier.",
            Category.SYSTEM, Category.ID),

    /** Direct-ID format identifier. */
    DIRECT_ID ("Direct-ID format identifier.",
            Category.SYSTEM, Category.ID),

    /** Direct-ID format identifier. */
    GRAPHICAL_ID ("Grapchical identifier.",
            Category.SYSTEM, Category.ID),

    /** HTML identifier. */
    H5_ID ("HTML identifier.",
          Category.SYSTEM, Category.ID),

    /** Direct-ID format identifier. */
    SCR_X ("X coordinate on Screen.",
            Category.REGULAR, Category.POSITIONING),

    /** Direct-ID format identifier. */
    SCR_Y ("Y coordinate on Screen.",
            Category.REGULAR, Category.POSITIONING),

    /** Root element for Direct-ID format identifier. */
    DIRECT_ID_ROOT ("Root element for Direct-ID format identifier.",
            Category.SYSTEM, Category.ID),

    /** Constant that determines the screen context to seek a component. */
    CONTEXT("Constant that determines the screen context to seek a component.",
            Category.SYSTEM, Category.AUTOMATION),

    /** Constant identifier of basic screen component. */
    GROUPROLE("Constant identifier of ID Group role.",
            Category.SYSTEM, Category.AUTOMATION),

    /** The title of the component */
    TITLE ("The title of the component",
            Category.REGULAR, Category.INFO),

    /** The title of the component */
    SCROLL_PAGE_SIZE ("The size of scroll page",
             Category.REGULAR, Category.INFO),

    /** The min position of the scroller */
    SCROLL_MIN_POS ("The minimum position",
             Category.REGULAR, Category.INFO),

    /** The max position of the scroller */
    SCROLL_MAX_POS ("The maximum position",
             Category.REGULAR, Category.INFO),

    /** The scroller actual position */
    SCROLL_POSITION ("The scroller position",
          Category.REGULAR, Category.INFO),

    /** The minimize direction of a dockable view */
    MINIMIZE_DIRECTION ("The dockable view's minimize direction",
          Category.REGULAR, Category.INFO)
    ;

    /**
     * This is a set of property categories in which this property is a member.
     */
    public final HashSet<Category> category;
    /**
     * Textual description of the property semantics.
     */
    public final String description;

    private Property(String description, Category...category) {
        this.description = description;
        this.category =
            new HashSet<Category>(Arrays.asList(category));
    }


    // =====================================================================
    // Property Filters convenience factory methods
    // =====================================================================
    /**
     * A factory method that creates a property-value filter to checks that:
     * <li> the property value as <code>String</code>
     * <b>is equal to</b> the provided text
     * <br><br>
     * <b>NOTES:</b>
     * <li> the parameter is transformed to {@link String}
     * <li> if the parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * @param textHolder - a string that must be exactly the same as the
     * property value.
     * @return an instance of property filter
     */
    public final IDPart stringEquals(Object textHolder) {
        return IDPart.from(this, MatchType.S_EQUALS.pack(textHolder));
    }
    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value as <code>String</code>
     * <b>starts with</b> the provided text
     * <br><br>
     * <b>NOTES:</b>
     * <li> the parameter is transformed to {@link String}
     * <li> if the parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * @param textHolder - a string that must be found exactly the same in the
     * start of the property value.
     * @return an instance of property filter
     */
    public final IDPart stringStarts(Object textHolder) {
        return IDPart.from(this, MatchType.S_STARTS.pack(textHolder));
    }
    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value as <code>String</code>
     * <b>ends with</b> the provided text
     * <br><br>
     * <b>NOTES:</b>
     * <li> the parameter is transformed to {@link String}
     * <li> if the parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * @param textHolder - a string that must be found exactly the same in the
     * end of the property value.
     * @return an instance of property filter
     */
    public final IDPart stringEnds(Object textHolder) {
        return IDPart.from(this, MatchType.S_ENDS.pack(textHolder));
    }
    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value as <code>String</code>
     * <b>contains</b> the provided text
     * <br><br>
     * <b>NOTES:</b>
     * <li> the parameter is transformed to {@link String}
     * <li> if the parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * @param textHolder - a string that must be found exactly the same in the
     * property value.
     * @return an instance of property filter
     */
    public final IDPart stringContains(Object textHolder) {
        return IDPart.from(this, MatchType.S_CONTAINS.pack(textHolder));
    }
    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value as <code>String</code>
     * <b>matches the Regular Expression</b> provided as text
     * <br><br>
     * <b>NOTES:</b>
     * <li> the parameter is transformed to {@link String}
     * <li> if the parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * @param regExHolder - a regular expression pattern, which the property
     * value must match.
     * @return an instance of property filter
     * @see Pattern
     */
    public final IDPart stringMatches(Object regExHolder) {
        return IDPart.from(this, MatchType.S_MATCHES.pack(regExHolder));
    }

    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value converted to <code>Number</code>
     * <b>is equal to</b> the provided comparison number
     * <br><br>
     * <b>NOTES:</b>
     * <li> the parameter is transformed to {@link String}
     * <li> if the parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * <li> in regards of the parameter's string representation - the
     * comparison with the property value is made in {@link Long} or
     * {@link Double} number format.
     * (e.g. "5.0" in {@link Double}, "5" in {@link Long})
     * <br><br>
     * @param numberHolder - a number to be matched by the property value.
     * @return an instance of property filter
     */
    public final IDPart numberEquals(Object numberHolder) {
        return numberEquals(numberHolder, 0.0);
    }
    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value converted to <code>Number</code> <b>is
     * into the interval determined by <tt>lowerLimitHolder</tt> and
     * <tt>upperLimitHolder</tt> inclusive</b>.
     * <br><br>
     * <b>NOTES:</b>
     * <li> all parameters are transformed to {@link String}
     * <li> if a parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * <li> in regards of the parameter's string representation - the
     * comparison with the property value is made in {@link Long} or
     * {@link Double} number format.
     * (e.g. "5.0" in {@link Double}, "5" in {@link Long})
     * <li> a boundary number of <b>null</b> is ignored (not used)
     * <br><br>
     * @param lowerLimitHolder - a number for lower inclusive limit to be
     * matched by the property value.
     * @param upperLimitHolder - a number for upper inclusive limit to be
     * matched by the property value.
     * @return an instance of property filter
     */
    public final IDPart numberFitsInclusive(
            Object lowerLimitHolder, Object upperLimitHolder) {
        return IDPart.from(this,
                MatchType.N_FITS_INCL.pack(lowerLimitHolder, upperLimitHolder));
    }
    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value converted to <code>Number</code> <b>is
     * into the interval determined by <tt>lowerLimitHolder</tt> and
     * <tt>upperLimitHolder</tt> exclusive</b>.
     * <br><br>
     * <b>NOTES:</b>
     * <li> all parameters are transformed to {@link String}
     * <li> if a parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * <li> in regards of the parameter's string representation - the
     * comparison with the property value is made in {@link Long} or
     * {@link Double} number format.
     * (e.g. "5.0" in {@link Double}, "5" in {@link Long})
     * <li> a boundary number of <b>null</b> is ignored (not used)
     * <br><br>
     * @param lowerLimitHolder - a number for lower exclusive limit to be
     * matched by the property value.
     * @param upperLimitHolder - a number for upper exclusive limit to be
     * matched by the property value.
     * @return an instance of property filter
     */
    public final IDPart numberFitsExclusive(
            Object lowerLimitHolder, Object upperLimitHolder) {
        return IDPart.from(this,
                MatchType.N_FITS_EXCL.pack(lowerLimitHolder, upperLimitHolder));
    }
    /**
     * This factory method creates a property-value filter to checks that:
     * <li> the property value converted to <code>Number</code> <b>is
     * near a given number within a given distance</b>
     * <br><br>
     * <b>NOTES:</b>
     * <li> all parameters are transformed to {@link String}
     * <li> if a parameter is of type {@link AbstractProperty} - the
     * contained value is first unwrapped
     * <li> in regards of the parameter's string representation - the
     * comparison with the property value is made in {@link Long} or
     * {@link Double} number format.
     * (e.g. "5.0" in {@link Double}, "5" in {@link Long})
     * <br><br>
     * @param numberHolder - a number that is the center of the match range.
     * @param deviationHolder - a number setting the maximum distance from the
     * central number where the property value must fit.
     * @return an instance of property filter
     */
    public final IDPart numberEquals(
            Object numberHolder, Object deviationHolder) {
        return IDPart.from(this,
                MatchType.N_EQUALS.pack(numberHolder, deviationHolder));
    }

    // =====================================================================
    // Property-bound utilities
    // =====================================================================
    /**
     * This map contains a list of properties which are members of a given
     * property category. This map is intended for use with the component
     * dumping tool.<br>
     * <b>NOTE:</b> One property could be a member of more than one category.
     */
    public static final HashMap<Category, List<Property>> CAT2PROP =
        new HashMap<Category, List<Property>>();

    static {
        for (Property prop : Property.values()) {
            for (Category cat : prop.category) {
                if (!CAT2PROP.containsKey(cat)) {
                    CAT2PROP.put(cat, new ArrayList<Property>());
                }
                CAT2PROP.get(cat).add(prop);
            }
        }
    }

    /**
     * This class contains converter functions intended to equalize string
     * representation of list values. The implemented methods concern the
     * interpretation of the "selected item" value for tree and grid components.
     * As far as they are lists - the included methods offer uniform way for
     * their conversion to/from string.
     *
     * @author dkozhuharov
     */
    public static class Convert {
        // Tree related converters
        public static final String TREE_PATH_DELIMITER = ">";
        public static String treePathToString(List<String> pathItems) {
            StringBuilder sb = new StringBuilder();
            if (pathItems != null) {
                for (String item : pathItems) {
                    sb.append(item).append(TREE_PATH_DELIMITER);
                }
                sb.deleteCharAt(sb.length() - 1);
            }
            return sb.toString();
        }
        public static List<String> stringToTreePath(String pathItemsString) {
            if (pathItemsString == null) {
                pathItemsString = "";
            }
            return new ArrayList<String>(
                    Arrays.asList(
                            pathItemsString.split("\\" + TREE_PATH_DELIMITER)
                    )
            );
        }

        // Grid related converters
        public static final String ROW_ITEM_DELIMITER = "|";
        public static String rowItemsToString(List<String> rowItems) {
            StringBuilder sb = new StringBuilder();
            if (rowItems != null) {
                for (String item : rowItems) {
                    sb.append(item).append(ROW_ITEM_DELIMITER);
                }
                sb.deleteCharAt(sb.length() - 1);
            }
            return sb.toString();
        }
        public static List<String> stringToRowItems(String rowItemsString) {
            if (rowItemsString == null) {
                rowItemsString = "";
            }
            return new ArrayList<String>(
                    Arrays.asList(
                            rowItemsString.split("\\" + ROW_ITEM_DELIMITER)
                    )
            );
        }
    }
}
