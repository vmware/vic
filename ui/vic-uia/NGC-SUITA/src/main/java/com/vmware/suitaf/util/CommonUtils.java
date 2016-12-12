package com.vmware.suitaf.util;

import java.lang.reflect.Method;
import java.util.HashMap;
import java.util.HashSet;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import com.vmware.hsua.common.datamodel.AbstractProperty;
import com.vmware.hsua.common.datamodel.PropertyBox;

/**
 * This class is a collection of commonly used tools.
 *
 * @author dkozhuharov
 */
public abstract class CommonUtils {
    // =================================================================
    // Value comparison instrumental methods
    // =================================================================
    /**
     * This method compares two objects and tries to determine their equality
     * the "smart" way. The two objects are considered equal in one of the
     * cases:
     * <li> <b>valie1</b> == null and <b>value2</b> == null
     * <li> <b>valie1.equals(value2)</b>
     * <li> <b>value1.equals( value1.getClass().valueOf(value2.toString()))</b>
     * <li> <b>value2.equals( value2.getClass().valueOf(value1.toString()))</b>
     * <br>
     * <b>NOTE:</b> if any of the parameters is an instance of
     * {@link AbstractProperty} - it is replaced by the wrapped property
     * value through calling <b>value?.get()</b><br>
     * <br>
     * This method is very helpful when comparing the values of some simple
     * types like {@link Boolean}, {@link Integer} or {@link Enum} to their
     * converted to string version. This will be often the case when comparing
     * the property values of UI components.<br>
     * Other advantage is that it could unwrap property field from the
     * {@link PropertyBox} containers.
     * <br>
     * @param value1 - the first value to be compared
     * @param value2 - the second value to be comapred
     * @return <b>true</b>
     */
    public static boolean smartEqual(
            Object value1, Object value2) {
        // If any one is null --> equal only if both are null
        if (value1 == null || value2 == null)
            return value1 == value2;

        // Apply value unwrapping for properties
        if (value1 instanceof AbstractProperty<?>)
            value1 = ((AbstractProperty<?>) value1).get();

        // Apply value unwrapping for properties
        if (value2 instanceof AbstractProperty<?>)
            value2 = ((AbstractProperty<?>) value2).get();

        // If simple equals() brings true --> equal
        if (value1.equals(value2))
            return true;

        // If valueOf transformed equals() brings true --> equal
        if (equalValueOf(value1, value2))
            return true;

        // If valueOf transformed equals() brings true --> equal
        if (equalValueOf(value2, value1))
            return true;

        // Values are --> different
        return false;
    }
    private static boolean equalValueOf(Object typedValue, Object untypedValue){
        Object converted = null;
        try {
            converted = valueOf(typedValue.getClass(), untypedValue);
        } catch (Exception e) {  }
        return typedValue.equals(untypedValue) || typedValue.equals(converted);
    }

    private static final HashMap<Class<?>, Class<?>> WRAPPERS =
        new HashMap<Class<?>, Class<?>>();
    static {
        WRAPPERS.put(java.lang.Boolean.TYPE, java.lang.Boolean.class);
        WRAPPERS.put(java.lang.Character.TYPE, java.lang.Character.class);
        WRAPPERS.put(java.lang.Byte.TYPE, java.lang.Byte.class);
        WRAPPERS.put(java.lang.Short.TYPE, java.lang.Short.class);
        WRAPPERS.put(java.lang.Integer.TYPE, java.lang.Integer.class);
        WRAPPERS.put(java.lang.Long.TYPE, java.lang.Long.class);
        WRAPPERS.put(java.lang.Float.TYPE, java.lang.Float.class);
        WRAPPERS.put(java.lang.Double.TYPE, java.lang.Double.class);
    }

    @SuppressWarnings("unchecked")
    public static <T> T valueOf(Class<T> targetClass, Object untypedValue)
    throws Exception {
        if (targetClass.equals(String.class)) {
            return (T) untypedValue.toString();
        }
        else if (targetClass.isPrimitive()) {
            String name = targetClass.getSimpleName();
            name = "parse"
                + name.substring(0,1).toUpperCase() + name.substring(1);
            Method method =
                WRAPPERS.get(targetClass).getDeclaredMethod(name, String.class);
            return (T) method.invoke(null, untypedValue.toString());
        }
        else {
            String methodName = "valueOf";
            Method method =
                targetClass.getDeclaredMethod(methodName, String.class);
            return (T) method.invoke(null, untypedValue.toString());
        }
    }

    // ===============================================================

    public static void sleep(long milliseconds) {
        try {
            Thread.sleep(milliseconds);
        } catch (InterruptedException e) {  }
    }

    private static HashSet<Character> GOOD_CHARS = new HashSet<Character>();
    static {
        char[] chArr =
            (" " +
            "abcdefghijklmnopqrstuvwxyz" +
            "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
            "0123456789").toCharArray();
        for (Character ch : chArr)
            GOOD_CHARS.add(ch);
    }

    public static String stripNonAlphaNumericChars(String s) {
        String result = "";
        for ( int i = 0; i < s.length(); i++ ) {
            if ( GOOD_CHARS.contains(s.charAt(i)) )
                result += s.charAt(i);
        }
        return result;
    }

    /**
     * Method that generates relatively unique random string with requested
     * length
     * @param length - of the returned random string
     * @return the random string
     */
    public static String getRandomString(int length) {
        String myRandom = "";
        while (myRandom.length() < length){
            UUID uuid = UUID.randomUUID();
            myRandom = myRandom + uuid.toString();
        }
        return myRandom.substring(0, length);
    }

    /**
     * This method is to be used for text normalization of human readable texts.
     * It allows matching of semantically equivalent texts that differ just by
     * punctuation, letter-case or white spaces.
     * <br><br>
     * The method is UTF-8 compatible.
     *
     * @param stringHolder - an object whose {@code  toString()} method is
     * called to derive the string content.
     *
     * @return Normalized version of the string
     */
    public static String normalize(Object stringHolder) {
        if (stringHolder==null)
            return null;
        else
            return stringHolder.toString()
                .replaceAll("([^\\d\\p{javaUpperCase}\\p{javaLowerCase}])+"," ")
                .trim()
                .toLowerCase();
    }

    /**
     * This method gets the strings of the two parameters normalizes them using
     * call to the {@link StringHelper#normalize(Object)} method and compares
     * if they are equal. If any of them is <b>null</b> the result is always
     * <b>false</b>.
     * @param firstStringHolder - first string source for the compare
     * @param secondStringHolder - second string source for the compare
     * @return <b>true</b> only if both strings are not null and after the
     * normalization they boil down to equal strings.
     */
    public static boolean equalsNormalized(
            Object firstStringHolder, Object secondStringHolder) {
        if (firstStringHolder==null || secondStringHolder==null)
            return false;
        else
            return normalize(firstStringHolder).equals(
                    normalize(secondStringHolder));
    }
    /**
     * Convert unicoded list of symbols to string.
     * @param   unicodeString
     * @return
     */
    public static String convertUnicodeToString(String unicodeString) {
        StringBuffer result = new StringBuffer();
        Matcher matcher = Pattern.compile("\\\\u((?i)[0-9a-f]{4})").matcher(unicodeString);
        while (matcher.find()) {
            int codepoint = Integer.valueOf(matcher.group(1), 16);
            result.append((char) codepoint);
        }
        return result.toString();
    }
}
