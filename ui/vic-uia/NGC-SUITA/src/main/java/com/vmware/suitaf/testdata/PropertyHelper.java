package com.vmware.suitaf.testdata;

import java.io.File;
import java.io.FileInputStream;
import java.lang.reflect.Array;
import java.lang.reflect.Field;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.TreeMap;

import com.vmware.hsua.common.datamodel.AbstractProperty;
import com.vmware.hsua.common.datamodel.PropertyBox;
import com.vmware.hsua.common.datamodel.AbstractProperty.Load;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.suitaf.util.Logger;

/**
 * This class contains the tool-set for loading of configuration values
 * from property files directly into class fields. A single class instance wraps
 * one property file.<br>
 * For wrapping - each property file is given with a relative file name. The
 * root folder for all property files is by default taken from the "user.dir"
 * system property (see {@link System#getProperties()}). It could be changed
 * with the setter method {@link #setPropertyRootFolder(File)}.<br>
 * <br>
 * Property files handled by this class must meet the following conditions:<br>
 * <ul>
 * <li> Each key in the property file must match the canonical name of
 * the class field for which the corresponding value is dedicated.
 * <li> If more than one value must be provided for a field - the keys will be
 * in the form: {@code <CanonicalFieldName>.<OrderingNumber>}
 * <li> The target fields must be of any of the following types:<br>
 *   <ul>
 *   <li> String
 *   <li> All the primitive types
 *   <li> All wrappers types: Boolean, Character, Byte, Short, Integer, Long,
 *   Float, Double
 *   <li> Enum-s
 *   <li> Any custom-defined type that implements a static function
 *   <tt>valueOf(String)</tt> for which the following is always true:
 *   <tt>val.equals(val.valueOf(val.toString()))</tt>
 *   </ul>
 * <li> The target fields could be arrays of the above listed types.
 * <li> The target fields could be {@link AbstractProperty}s of the above
 * listed types. Any {@link AbstractProperty} field that is will receive value
 * from a property-file must be annotated with the {@link Load} annotation.
 * <li> The target fields must <b>not</b> be <b><tt>final</tt></b>. This means
 * also that they could not be declared in an interface.
 * <li> The target fields can be <b><tt>private</tt></b> or
 * <b><tt>protected</tt></b>.
 * </ul>
 * <br>
 *
 * @author dkozhuharov
 */
public final class PropertyHelper {
    private static File PROPERTY_ROOT_FOLDER = getPropertyRootFolderDefault();
    /**
     * This method sets the root folder for property-files to the default one
     * provided by method {@link #getPropertyRootFolderDefault()}.
     */
    public static final void setPropertyRootFolderDefault() {
        setPropertyRootFolder(getPropertyRootFolderDefault());
    }
    /**
     * This method sets the root folder for property-files to the one provided
     * as input parameter.
     * @param propertyRootFolder - the new root folder for property-files
     * @throws RuntimeException if the provided {@link File} is not a folder.
     */
    public static final void setPropertyRootFolder(File propertyRootFolder) {
        PROPERTY_ROOT_FOLDER = propertyRootFolder;
        if (!PROPERTY_ROOT_FOLDER.isDirectory()) {
            throw new RuntimeException(
                    "Property root folder not found: " + PROPERTY_ROOT_FOLDER);
        }
    }
    /**
     * This method retrieves the current root folder for property-files.
     * @return current root folder for property-files
     */
    public static final File getPropertyRootFolder() {
        return PROPERTY_ROOT_FOLDER;
    }
    /**
     * This method retrieves the default root folder for property-files. It is
     * based on the "user.dir" system property of Java.
     * @return default root folder for property-files
     */
    public static final File getPropertyRootFolderDefault() {
        return new File(System.getProperty("user.dir"));
    }


    static final File getPropertyFile(String propertyFileName) {
        return new File(getPropertyRootFolder(), propertyFileName);
    }
    static final HashMap<File, Properties> PropertyCache =
        new HashMap<File, Properties>();
    static final Object PropertyCacheLock =
        new Object();
    static final Properties getProperties(File propertyFile) {
        Properties properties = null;

        synchronized (PropertyCacheLock) {
            if (PropertyCache.containsKey(propertyFile)) {
                properties = PropertyCache.get(propertyFile);
            }
            else {
                properties = new Properties();
                try {
                    properties.load(new FileInputStream(propertyFile));
                } catch (Exception e) {
                    throw new RuntimeException(
                            "The property file could not be loaded", e);
                }
                PropertyCache.put(propertyFile, properties);
            }
        }

        return properties;
    }

    @SuppressWarnings("serial")
    private static final class PropertyMap
    extends HashMap<Class<?>, Map<Field, TreeMap<Integer,String>>> {};

    private static final PropertyMap parseProperties(
            String propertyFileName, Class<?> filterClass) {

        boolean hasFailingFields = false;
        Properties prop = getProperties(getPropertyFile(propertyFileName));

        PropertyMap parsedMap = new PropertyMap();
        HashMap<String, Class<?>> className2Inst =
            new HashMap<String, Class<?>>();

        // If filter class is given
        if (filterClass != null) {
            className2Inst.put(filterClass.getCanonicalName(), filterClass);
        }

        for (Object keyObject : prop.keySet()) {
            String keyText = keyObject.toString();
            int fldPoint = keyText.lastIndexOf(".");
            if (fldPoint < 0) {
                continue;
            }

            // Extract (if exists) the trailing numbers as they are to be used
            // for value sorting
            Integer valueKey = 0;
            if (keyText.substring(fldPoint + 1).matches("[0-9]+")) {
                valueKey = Integer.valueOf(keyText.substring(fldPoint + 1));
                keyText = keyText.substring(0, fldPoint);
                fldPoint = keyText.lastIndexOf(".");
                if (fldPoint < 0) {
                    continue;
                }
            }

            String className = keyText.substring(0, fldPoint);
            String fieldName = keyText.substring(fldPoint + 1);
            String value = prop.get(keyObject).toString();

            // Try to find the class and the field for which the property value
            // is addressed
            try {
                if (filterClass == null) {
                    if (!className2Inst.containsKey(className)) {
                        Class<?> classInst = null;
                        String tmpName = className;
                        // Conversion of "canonical" into "normal" class name
                        while (classInst == null && tmpName.contains(".")) {
                            try {
                                classInst = Class.forName(tmpName);
                            } catch (Exception e) {
                                int i = tmpName.lastIndexOf(".");
                                tmpName = tmpName.substring(0, i) + "$" +
                                tmpName.substring(i + 1);
                            }
                        }
                        if (classInst == null) {
                            throw new ClassNotFoundException(className);
                        }
                        else {
                            className2Inst.put(className, classInst);
                        }
                    }
                }

                // Finding the class
                Class<?> classInst = className2Inst.get(className);
                if (classInst == null) {
                    continue;
                }

                // Finding the field - if not exiting - throw exception
                Field fieldInst = classInst.getDeclaredField(fieldName);

                // add the found class instance
                if (!parsedMap.containsKey(classInst)) {
                    parsedMap.put(classInst,
                            new HashMap<Field, TreeMap<Integer,String>>());
                }

                // add the found field instance
                if (!parsedMap.get(classInst).containsKey(fieldInst)) {
                    parsedMap.get(classInst).put(
                            fieldInst, new TreeMap<Integer,String>());
                }

                // add the property-value with the key for value sorting
                parsedMap.get(classInst).get(fieldInst).put(valueKey, value);
            } catch (Exception e) {
                Logger.error("Property parsing failed for: " +
                        keyText + "=" + value, e);
                hasFailingFields = true;
            }

        }

        if (hasFailingFields) {
            RuntimeException e = new RuntimeException("Some fields parsing" +
                    " failed for property file: " + propertyFileName);
            Logger.error(e.getMessage());
            throw e;
        }

        return parsedMap;
    }

    private static final void setFields(String propertyFileName,
            Class<?> targetClass, Object targetInstance) {
        boolean hasFailingFields = false;
        if (targetClass == null && targetInstance != null) {
            targetClass = targetInstance.getClass();
        }
        PropertyMap parsedMap = parseProperties(propertyFileName, targetClass);
        for (Class<?> currentClass : parsedMap.keySet()) {
            for (Field fld : parsedMap.get(currentClass).keySet()) {
                Collection<String> values =
                    parsedMap.get(currentClass).get(fld).values();
                Class<?> type = fld.getType();
                try {
                    if (AbstractProperty.class.isAssignableFrom(type)) {
                        ((AbstractProperty<?>)
                                fld.get(targetInstance)).load(values);
                    }
                    else if (type.isArray()) {
                        type = type.getComponentType();
                        Object arr = Array.newInstance(type, values.size());
                        int i=0;
                        for (String value : values) {
                            Array.set(arr, i++,
                                    CommonUtils.valueOf(type, value));
                        }
                        fld.set(targetInstance, arr);
                    }
                    else {
                        String value = values.iterator().next();
                        Object fldValue = CommonUtils.valueOf(type, value);

                        if (fld.isAccessible()) {
                            fld.set(targetInstance, fldValue);
                        }
                        else {
                            fld.setAccessible(true);
                            fld.set(targetInstance, fldValue);
                            fld.setAccessible(false);
                        }
                    }
                } catch (Exception e) {
                    Logger.error("Property setting failed for: " +
                            fld + "=" + values, e);
                    hasFailingFields = true;
                }
            }
        }

        if (hasFailingFields) {
            RuntimeException e = new RuntimeException("Some values loading" +
                    " failed for property file: " + propertyFileName);
            Logger.error(e.getMessage());
            throw e;
        }
    }

    //========================================================================
    // PropertyLoader instance declarations (non-static fields & methods)
    //========================================================================
    private final String relativeFileName;
    private PropertyHelper(String relativeFileName) {
        this.relativeFileName = relativeFileName;
    }

    /**
     * This is a factory method that creates {@link PropertyHelper} instances
     * wrapping one property file. The property file is given as a parameter
     * with its relative file name. The name must be relative to the common
     * property-file root folder (see {@link #setPropertyRootFolder(File)}).
     * @param relativeFileName - relative file name of a property file
     * @return a {@link PropertyHelper} instances wrapping the property file.
     */
    public static final PropertyHelper get(String relativeFileName) {
        return new PropertyHelper(relativeFileName);
    }

    /**
     * This is a factory method that creates {@link PropertyBox} instances.
     * Each instance is populated with the property values found in the
     * wrapped property file.
     * @param <T> - this type parameter determines the type of the PropertyBox
     * to be created.
     * @param pboxClass - the class of the property box to be created
     * @return an instance of a property box from the requested type
     */
    public final <T extends PropertyBox> T loadPBox(Class<T> pboxClass) {
        T pbox = null;
        try {
            pbox = pboxClass.newInstance();
        } catch (Exception e) {
            throw new RuntimeException("PBox instantiation failed for: " +
                    relativeFileName + "/" + pboxClass, e);
        }
        setFields(relativeFileName, pboxClass, pbox);
        return pbox;
    }
    /**
     * This method loads all values from the wrapped property file into their
     * destination fields.<br>
     * <b>NOTE:</b> All the fields described in the property file must be
     * static. Otherwise exception would be thrown.<br>
     * <b>NOTE:</b> The target fields must not be final.<br>
     */
    public final void loadFields() {
        setFields(relativeFileName, null, null);
    }
    /**
     * This method loads the values from the wrapped property file into their
     * destination fields. Only the fields from the requested target class
     * are processed. All other fields are skipped.<br>
     * <b>NOTE:</b> All the fields described in the property file, which are
     * from the target class, must be static. Otherwise exception would be
     * thrown.<br>
     * <b>NOTE:</b> The target fields must not be final.<br>
     * @param targetClass - the class for which the loading would be processed.
     */
    public final void loadFields(Class<?> targetClass) {
        setFields(relativeFileName, targetClass, null);
    }
    /**
     * This method loads the values from the wrapped property file into their
     * destination fields. Only the fields from the class of the given instance
     * are processed. All other fields are skipped.<br>
     * <b>NOTE:</b> The fields described in the property file, which are
     * from the target instance class could be both static or dynamic.<br>
     * <b>NOTE:</b> The target fields must not be final.<br>
     * @param targetInstance - the instance to which the values will be loaded.
     */
    public final void loadFields(Object targetInstance) {
        setFields(relativeFileName, null, targetInstance);
    }
}
