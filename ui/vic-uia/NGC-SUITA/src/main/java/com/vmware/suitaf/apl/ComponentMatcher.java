package com.vmware.suitaf.apl;

/**
 * The {@link ComponentMatcher} interface will be used by all {@link
 * AutomationPlatformLink} implementors to mark their specific component
 * matcher implementations. The factory method {@link AutomationPlatformLink
 * #getComponentMatcher(IDGroup...)} will implement matcher instanciations.
 * <br><br>
 * The role of the component matcher is to allow identification and
 * operations with one main screen component. It optionally allows
 * identification of other related screen components which are needed for
 * improvement of the identification or the operations with the main
 * component.
 * <br><br>
 * To fulfill these roles the Component Matcher steps on a list of
 * {@link IDGroup} components. An ID group holds the information needed for
 * the identification of one screen component. It is in the form of
 * Property-Value pairs that are used to match against the properties of
 * the screen components. Besides the {@link Category#REGULAR}
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
 * Property-Value pairs are packed in {@link IDPart} instances. The
 * available properties are enlisted in the enumeration {@link Property}.
 * The property values are stored in the form of String representations.
 * <br><br>
 * @author dkozhuharov
 */
public interface ComponentMatcher {
    /**
     * This methods implementation will help logging of failure states and
     * debugging. It must provide meaningful, short and unique representation
     * of the native identifier.
     * <br><br>
     * @return one line text representation
     */
    public abstract String getLogForm();
}