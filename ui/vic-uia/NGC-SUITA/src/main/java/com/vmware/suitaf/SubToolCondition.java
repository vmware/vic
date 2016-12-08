package com.vmware.suitaf;

import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;
import com.vmware.suitaf.apl.SpecialStateHandler;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.suitaf.util.Condition;

/**
 * This is the sub-tool {@link SubToolCondition} with the following specs:
 * <li> <b>Function Type:</b> COMPONENT FACTORY
 * <li> <b>Description:</b> Creation of {@link Condition} instances to be used
 * for checking, awaiting or assertion of specific conditions of the UI or
 * of the test application.
 * <li> <b>Based on SubTool:</b>
 * {@link SubToolComponent#exists(Object, IDPart...)}
 */
public class SubToolCondition extends BaseSubTool {
    public SubToolCondition(UIAutomationTool ui) {
        super(ui);
    }
    // ======================================================================
    // Factory methods that generate simple UI-related conditions
    // ======================================================================
    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> IS UI component with the specified id FOUND on screen
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition isFound(final Object id, final IDPart...propertyFilters) {
        return new Condition(
                "IS_Found<" + id + arrayToString(propertyFilters) + ">") {
            @Override
            protected boolean checkImpl() {
                return ui.component.exists(id, propertyFilters);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> IS the special state of the specified {@link SpecialStateHandler}
     * FOUND on screen
     * @param ssh - a {@link SpecialStateHandler} instance
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition isFound(final SpecialStateHandler ssh) {
        return new Condition("IS_Found<" + ssh + ">" ) {
            @Override
            protected boolean checkImpl() {
                return ssh.stateRecognize();
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> is UI component with the specified id NOT FOUND on screen
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition notFound(final Object id, final IDPart...propertyFilters) {
        return new Condition(
                "NOT_Found<" + id + arrayToString(propertyFilters) + ">") {
            @Override
            protected boolean checkImpl() {
                return !ui.component.exists(id, propertyFilters);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> is the special state of the specified {@link SpecialStateHandler}
     * NOT FOUND on screen
     * @param ssh - a {@link SpecialStateHandler} instance
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition notFound(final SpecialStateHandler ssh) {
        return new Condition("NOT_Found<" + ssh + ">" ) {
            @Override
            protected boolean checkImpl() {
                return !ssh.stateRecognize();
            }
        };
    }

    // ======================================================================
    // Factory methods that generate simple value-related conditions
    // ======================================================================
    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> IS the specified "actual" object SAME as the specified "expected"
     * object. The comparison is done using the method {@link CommonUtils
     * #smartEqual(Object, Object)}, which makes it very flexible and
     * comfortable for use.
     * @param actual - the actual value at the time of comparison
     * @param expected - the expected value considered earlier
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition isSame(final Object actual, final Object expected) {
        return new Condition(
                "IS_Same<" + actual + ", " + expected + ">") {
            @Override
            protected boolean checkImpl() {
                return CommonUtils.smartEqual(actual, expected);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> is the specified "actual" object NOT SAME as the specified "expected"
     * object. The comparison is done using the method {@link CommonUtils
     * #smartEqual(Object, Object)}, which makes it very flexible and
     * comfortable for use.
     * @param actual - the actual value at the time of comparison
     * @param expected - the expected value considered earlier
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition notSame(final Object actual, final Object expected) {
        return new Condition(
                "NOT_Same<" + actual + ", " + expected + ">") {
            @Override
            protected boolean checkImpl() {
                return !CommonUtils.smartEqual(actual, expected);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> IS the specified "value" object TRUE. The comparison is done using
     * the method {@link CommonUtils#smartEqual(Object, Object)}, which makes
     * it very flexible and comfortable for use.
     * @param value - the value to be evaluated
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition isTrue(final Object value) {
        return new Condition(
                "IS_True<" + value + ">") {
            @Override
            protected boolean checkImpl() {
                return CommonUtils.smartEqual(value, true);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> is the specified "value" object NOT TRUE. The comparison is done
     * using the method {@link CommonUtils#smartEqual(Object, Object)}, which
     * makes it very flexible and comfortable for use.
     * @param value - the value to be evaluated
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition notTrue(final Object value) {
        return new Condition(
                "NOT_True<" + value + ">") {
            @Override
            protected boolean checkImpl() {
                return !CommonUtils.smartEqual(value, true);
            }
        };
    }

    // ======================================================================
    // Factory methods that constructs compound conditions from several simple
    // ======================================================================
    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> ALL inner {@link Condition}s from a given list must be TRUE.
     * <br><br>
     * <b>NOTE:</b> The principle of short boolean evaluation is applied.
     * Inner conditions are evaluated sequentially until a <b>false</b>
     * condition is found or all appear to be <b>true</b>. This could lead to
     * evaluation of just part of the inner condition.
     * <br><br>
     * @param conditionList - list of inner conditions to be evaluated.
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition allTrue(final Condition...conditionList) {
        if (conditionList.length == 0) {
            throw new RuntimeException("At least one condition must be given.");
        }
        return new Condition("ALL_True") {
            @Override
            protected boolean checkImpl() {
                for (Condition cnd : conditionList) {
                    if (!cnd.estimate())  return false;
                }
                return true;
            }
            @Override
            public String toString() {
                return super.toString() + arrayToString(conditionList);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> there must be ANY of the inner {@link Condition}s that is TRUE.
     * <br><br>
     * <b>NOTE:</b> The principle of short boolean evaluation is applied.
     * Inner conditions are evaluated sequentially until a <b>true</b>
     * condition is found or all appear to be <b>false</b>. This could lead to
     * evaluation of just part of the inner condition.
     * <br><br>
     * @param conditionList - list of inner conditions to be evaluated.
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition anyTrue(final Condition...conditionList) {
        if (conditionList.length == 0) {
            throw new RuntimeException("At least one condition must be given.");
        }
        return new Condition("ANY_True") {
            @Override
            protected boolean checkImpl() {
                for (Condition cnd : conditionList) {
                    if (cnd.estimate())  return true;
                }
                return false;
            }
            @Override
            public String toString() {
                return super.toString() + arrayToString(conditionList);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> ALL inner {@link Condition}s from a given list must be FALSE.
     * <br><br>
     * <b>NOTE:</b> The principle of short boolean evaluation is applied.
     * Inner conditions are evaluated sequentially until a <b>true</b>
     * condition is found or all appear to be <b>false</b>. This could lead to
     * evaluation of just part of the inner condition.
     * <br><br>
     * @param conditionList - list of inner conditions to be evaluated.
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition allFalse(final Condition...conditionList) {
        if (conditionList.length == 0) {
            throw new RuntimeException("At least one condition must be given.");
        }
        return new Condition("ALL_False") {
            @Override
            protected boolean checkImpl() {
                for (Condition cnd : conditionList) {
                    if (cnd.estimate())  return false;
                }
                return true;
            }
            @Override
            public String toString() {
                return super.toString() + arrayToString(conditionList);
            }
        };
    }

    /**
     * Factory method that creates a {@link Condition} with success criterion:
     * <li> there must be ANY of the inner {@link Condition}s that is FALSE.
     * <br><br>
     * <b>NOTE:</b> The principle of short boolean evaluation is applied.
     * Inner conditions are evaluated sequentially until a <b>false</b>
     * condition is found or all appear to be <b>true</b>. This could lead to
     * evaluation of just part of the inner condition.
     * <br><br>
     * @param conditionList - list of inner conditions to be evaluated.
     * @return a {@link Condition} object capable of checking the criterion
     */
    public Condition anyFalse(final Condition...conditionList) {
        if (conditionList.length == 0) {
            throw new RuntimeException("At least one condition must be given.");
        }
        return new Condition("ANY_False") {
            @Override
            protected boolean checkImpl() {
                for (Condition cnd : conditionList) {
                    if (! cnd.estimate())  return true;
                }
                return false;
            }
            @Override
            public String toString() {
                return super.toString() + arrayToString(conditionList);
            }
        };
    }

    /**
     * Helper method that aids the logging of group of objects.
     * @param <T> - the generic type of the objects in the group
     * @param elements - vararg parameter to receive the objects
     * @return the string representation of the group
     */
    private <T extends Object> String arrayToString(T...elements) {
        StringBuilder description = new StringBuilder();
        description.append("(");
        for (int i=0; i<elements.length; i++) {
            Object element = elements[i];
            if (i > 0)  description.append(", ");
            description.append(element.toString());
        }
        description.append(")");
        return description.toString();
    }
}
