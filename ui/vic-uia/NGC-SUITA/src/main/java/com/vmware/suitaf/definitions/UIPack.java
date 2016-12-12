package com.vmware.suitaf.definitions;

import java.lang.annotation.Annotation;
import java.lang.annotation.Documented;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.reflect.Field;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

import com.vmware.suitaf.apl.IDGroup;

/**
 * This class defines the concept and the toolage for the <b>UI-Pack</b>s.
 *
 * @author dkozhuharov
 */
public class UIPack {
    // ===================================================================
    // Static support for UIPacks
    // ===================================================================
    private static final Map<Class<?>, UIPack> Class2UIPack =
        new HashMap<Class<?>, UIPack>();
    private static final Map<IDGroup, UIPack> ID2UIPack =
        new HashMap<IDGroup, UIPack>();

    /**
     * Static method for registering of an interface containing Component ID
     * Constants as a <b>UI-Pack</b>.
     * @param <E>
     * @param interfaceClass - a class instance representing the interface being
     * registered
     * @return always returns <b>null</b>.
     */
    public static <E> E newUIPack(Class<E> interfaceClass) {
        if (Class2UIPack.containsKey(interfaceClass)) {
            throw new RuntimeException("The UIPack was already registered." +
                    "Class: " + interfaceClass.getCanonicalName());
        }

        ArrayList<IDGroup> rootIDs = new ArrayList<IDGroup>();
        ArrayList<IDGroup> checkIDs = new ArrayList<IDGroup>();
        ArrayList<IDGroup> allIDs = new ArrayList<IDGroup>();

        for (Field idField : interfaceClass.getDeclaredFields()) {
            if (idField.getType().isAssignableFrom(IDGroup.class)) {
                try {
                    IDGroup idValue = (IDGroup) idField.get(null);
                    for (Annotation ann : idField.getAnnotations()) {
                        if (ann instanceof UIPackRoot) {
                            rootIDs.add(idValue);
                        }
                        if (ann instanceof UIPackCheck) {
                            checkIDs.add(idValue);
                        }
                    }
                    allIDs.add(idValue);
                } catch (Exception e) {
                    throw new RuntimeException("Unexpected access failure." +
                            " Field: " + idField, e);
                }
            }
        }

        // Create the new UIPack instance
        UIPack newUIPack =
            new UIPack(interfaceClass,
                    rootIDs.toArray(new IDGroup[0]),
                    checkIDs.toArray(new IDGroup[0]),
                    allIDs.toArray(new IDGroup[0]));

        // Store the new UIPack instance in the registers
        Class2UIPack.put(interfaceClass, newUIPack);
        for (IDGroup id : allIDs) {
            ID2UIPack.put(id, newUIPack);
        }

        // Return
        return null;
    }

    // ===================================================================
    // UIPack IDs role annotations
    // ===================================================================
    /**
     * This annotation is to be used in the {@link UIPack} definitions to
     * mark the component IDs that appears as root for all UIPack IDs.
     * <br><br>
     * @author dkozhuharov
     */
    @Retention(RetentionPolicy.RUNTIME)
    @Documented
    public static @interface UIPackRoot {}

    /**
     * This annotation is to be used in the {@link UIPack} definitions to
     * mark the component IDs to be used to check UIPack presence.
     * <br><br>
     * @author dkozhuharov
     */
    @Retention(RetentionPolicy.RUNTIME)
    @Documented
    public static @interface UIPackCheck {}

    // ===================================================================
    // UIPack fields
    // ===================================================================
    public final Class<?> definingClass;
    public final IDGroup[] rootIDs;
    public final IDGroup[] checkIDs;
    public final IDGroup[] allIDs;

    private UIPack(Class<?> definingClass,
            IDGroup[] rootIDs, IDGroup[] checkIDs, IDGroup[] allIDs) {
        this.definingClass = definingClass;
        this.rootIDs = rootIDs;
        this.checkIDs = checkIDs;
        this.allIDs = allIDs;
    }
}
