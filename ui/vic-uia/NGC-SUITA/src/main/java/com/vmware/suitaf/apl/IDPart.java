package com.vmware.suitaf.apl;

import java.util.Arrays;
import java.util.HashSet;

import com.vmware.hsua.common.datamodel.AbstractProperty;

/**
 * This class represents an ID part, which is a Property-Value pair. The
 * Property-Value pairs are used to match against the properties of the screen
 * components. Besides the {@link Category#REGULAR} Property-Value pairs used
 * for matching, there are also {@link Category#SYSTEM} pairs. The system pairs
 * are not matched directly to a component's property, but are used to control
 * the component identification and component actions.<br>
 * See {@link IDGroup} for further details.
 * <br><br>
 * <b>NOTE:</b> This class is designed to be immutable. Please do not change
 * this feature.
 * <br><br>
 * @author dkozhuharov
 */
public final class IDPart {

    /**
     * This <b>enum</b> represents the available property filtering operations.
     * The property filtering is part of the {@link IDGroup} abstract identifier
     * concept.<br>
     * All filtering operations are split on two sets - {@link String}-based
     * and {@link Number}-based.<br>
     * @author dkozhuharov
     *
     */
    static enum MatchType {
        S_EQUALS(false),
        S_STARTS(false),
        S_ENDS(false),
        S_CONTAINS(false),
        S_MATCHES(false),

        N_FITS_INCL(true),
        N_FITS_EXCL(true),
        N_EQUALS(true),
        ;

        public final boolean isNumeric;

        private MatchType(boolean isNumeric) {
            this.isNumeric = isNumeric;
        }


        private static final String PACK_SEPARATOR = "~I~";
        private static final String PACK_NULL_VALUE = "~NULL~";

        /**
         * Tool that strips the {@link AbstractProperty} wrapper if any.
         * @param property - a potential property-type value container
         * @return - the stripped value
         */
        private Object strip(Object property) {
            if (property instanceof AbstractProperty<?>) {
                return ((AbstractProperty<?>) property).get();
            }
            return property;
        }

        /**
         * This helper method works on a particular instance of the
         * {@link MatchType} <b>enum</b>. It packs the <b>enum</b> name and
         * the provided parameters as a single {@link String}. This packed
         * string plays the role of one property filtering command and is
         * to be used as value for an {@link IDPart} property filter.<br>
         * Each of the passed parameter objects is transformed using its
         * <tt>toString()</tt> method. If a parameter is of type
         * {@link AbstractProperty} - then it is first unpacked from the
         * property container and then transformed to string.
         *
         * @param args - a vararg parameter receiving open list of the
         * comparison command parameters.
         * @return a {@link String} representing the packed form of the
         * comparison command.
         */
        String pack(Object...args) {
            StringBuffer packed = new StringBuffer();
            packed.append(PACK_SEPARATOR).append(this.name());
            for (Object arg : args) {
                packed.append(PACK_SEPARATOR);
                if (arg != null) {
                    packed.append(strip(arg).toString());
                }
                else {
                    packed.append(PACK_NULL_VALUE);
                }
            }
            return packed.toString();
        }

        /**
         * This static helper function unpacks a packed property comparison
         * command and extracts the comparison command type.<br>
         * If the string format is not recognized then it is accepted for plain
         * {@link String} and the default comparison command of {@link
         * MatchType#S_EQUALS} is returned. This is done for backward
         * compatibility with the old property filtering mode.<br>
         * @param packed - a packed property comparison command
         * @return a {@link MatchType} value representing the type of the
         * encoded comparison command
         */
        static MatchType unpackType(String packed) {
            try {
                if (packed.startsWith(PACK_SEPARATOR)) {
                    // If package formatting is recognized
                    return Enum.valueOf( MatchType.class,
                            packed.split(PACK_SEPARATOR)[1] );
                }
            } catch (Exception e) {}
            // else - it's a raw text
            return S_EQUALS;
        }
        /**
         * This static helper function unpacks a packed property comparison
         * command and extracts the comparison command parameters. The
         * parameters are returned as an array of {@link String}s.
         * @param packed - a packed property comparison command
         * @return a {@link String} array containing the comparison command
         * parameters in {@link String} format or <b>null</b>.
         */
        static String[] unpackArgs(String packed) {
            if (packed.startsWith(PACK_SEPARATOR)) {
                // If package formatting is recognized
                String[] args = packed.split(PACK_SEPARATOR);
                String[] unpacked = new String[args.length - 2];
                for (int i=2; i<args.length; i++) {
                    int i_u = i-2;
                    if (PACK_NULL_VALUE.equals(args[i])) {
                        unpacked[i_u] = null;
                    }
                    else {
                        unpacked[i_u] = args[i];
                    }
                }
                return unpacked;
            }
            else {
                // else - it's a raw text
                return new String[] {packed};
            }
        }
    }


    /**
     * This field represents the property of this ID part
     */
    public final Property property;
    /**
     * This field holds a string for the ID part value
     */
    public final String value;

    private IDPart(Property property, String value) {
        this.property = property;
        this.value = value;
    }


    private static final HashSet<Property> CONSTANT_PROPERTIES =
        new HashSet<Property>( Arrays.asList(
                Property.CONTEXT, Property.GROUPROLE
        ) );

    /**
     * Factory method for instantiation of {@link IDPart}s. The method applies
     * different verifications of the input parameters.
     * @param property - identifier of {@link Property} for which this
     * {@link IDPart} instance will hold a matching value.
     * Must not be <b>null</b>
     * @param value    - value to be matched for the property.
     * @return an {@link IDPart} instance
     */
    public static IDPart from(Property property, String value) {
        if (CONSTANT_PROPERTIES.contains(property)) {
            throw new UnsupportedOperationException(
                    "Operation not supported for properties " +
                    CONSTANT_PROPERTIES + ". Use the predefined constants.");
        }
        if (property == null) {
            throw new IllegalArgumentException("Property must not be null");
        }
        return new IDPart(property, value);
    }

    /**
     * Factory method for instantiation of identifier {@link IDPart}. The method
     * checks if the input string starts with "//" or "./". If so an instance
     * with the {@link Property#X_PATH_ID} is created. Otherwise an instance
     * with the {@link Property#DIRECT_ID} is created.
     * @param stringID - value to be treated as an ID property.
     * @return an {@link IDPart} instance
     */
    public static IDPart id(String stringID) {
        if (stringID.startsWith("//") || stringID.startsWith("./")
                || stringID.startsWith("../") || stringID.equals(".")) {
            return new IDPart(Property.X_PATH_ID, stringID);
        }
        else {
            return new IDPart(Property.DIRECT_ID, stringID);
        }
    }

    /**
     * This is a constant-holding class. It contains {@link IDPart} instances
     * of the property {@link Property#CONTEXT}.
     * @author dkozhuharov
     */
    public static final class CONTEXT {
        public static final IDPart DESKTOP =
            new IDPart(Property.CONTEXT, "Desktop");

        public static final IDPart BROWSER_APP =
            new IDPart(Property.CONTEXT, "BrowserApplication");

        public static final IDPart BROWSER_WIN =
            new IDPart(Property.CONTEXT, "BrowserWindow");

        public static final IDPart WEB_APP =
            new IDPart(Property.CONTEXT, "WebApplication");

        private CONTEXT() {  }
    }

    /**
     * This is a constant-holding class. It contains {@link IDPart} instances
     * of the property {@link Property#GROUPROLE}.
     * @author dkozhuharov
     */
    public static final class GROUPROLE {
        /** GROUPROLE: Main component ID group */
        public static final IDPart MAIN =
            new IDPart(Property.GROUPROLE, "MAIN");

        /** GROUPROLE: Parent component ID group */
        public static final IDPart PARENT =
            new IDPart(Property.GROUPROLE, "PARENT");

        private GROUPROLE() {  }
    }

    @Override
    public String toString() {
        return property + "=" + value;
    }
}
