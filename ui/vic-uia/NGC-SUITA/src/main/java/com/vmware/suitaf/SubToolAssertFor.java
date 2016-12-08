/**
 *
 */
package com.vmware.suitaf;

import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.suitaf.util.Condition;
import com.vmware.suitaf.util.FailureCauseHolder;

/**
 * This is the sub-tool {@link SubToolAssertFor} with the following specs:
 * <li> <b>Function Type:</b> STATE ASSERTION
 * <li> <b>Description:</b> Contains commands that evaluate {@link Condition}s
 * and log the success if the condition was <b>true</b> or launch
 * the auditing procedure, log failure and fail the test if the condition
 * was <b>false</b>.
 * <li> <b>Based on SubTool:</b>
 * {@link SubToolCondition#isFound(Object, IDPart...)},
 * {@link SubToolCondition#notFound(Object, IDPart...)},
 * {@link SubToolCondition#isTrue(Object)},
 * {@link SubToolCondition#notTrue(Object)},
 * {@link SubToolCondition#isSame(Object, Object)},
 * {@link SubToolCondition#notSame(Object, Object)},
 * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
 */
public class SubToolAssertFor extends BaseSubTool {
    public SubToolAssertFor(UIAutomationTool ui) {
        super(ui);
    }

    /**
     * State Assertion Method that asserts a given fact by checking for the
     * condition:<br><b>
     * IS UI component with the specified id FOUND on screen
     * </b><br>
     * The condition is created using {@link SubToolCondition
     * #isFound(Object, IDPart...)} and is awaited to for "maxFactWait" ms.
     * The output of the assertion is handled through the auditing tool:
     * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
     * @param assertedFact - text describing the fact to be asserted
     * @param maxFactWait - maximum time to await the fact
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     */
    public void isFound(String assertedFact, long maxFactWait,
            Object id, IDPart...propertyFilters) {
        Condition c = ui.condition.isFound(id, propertyFilters);
        c.await(maxFactWait);
        condition(assertedFact, c);
    }

    /**
     * State Assertion Method that asserts a given fact by checking for the
     * condition:<br><b>
     * is UI component with the specified id NOT FOUND on screen
     * </b><br>
     * The condition is created using {@link SubToolCondition
     * #notFound(Object, IDPart...)} and is awaited to for "maxFactWait" ms.
     * The output of the assertion is handled through the auditing tool:
     * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
     * @param assertedFact - text describing the fact to be asserted
     * @param maxFactWait - maximum time to await the fact
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     */
    public void notFound(String assertedFact, long maxFactWait,
            Object id, IDPart...propertyFilters) {
        Condition c = ui.condition.notFound(id, propertyFilters);
        c.await(maxFactWait);
        condition(assertedFact, c);
    }

    /**
     * State Assertion Method that asserts a given fact by checking for the
     * condition:<br><b>
     * IS the specified "actual" object SAME as the specified "expected" object
     * </b><br> The comparison is done using the method
     * {@link CommonUtils#smartEqual(Object, Object)}, which makes it very
     * flexible and comfortable for use.<br>
     * The condition is created using {@link SubToolCondition
     * #isSame(Object, Object)} and is awaited to for "maxFactWait" ms.<br>
     * The output of the assertion is handled through the auditing tool:
     * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
     * @param assertedFact - text describing the fact to be asserted
     * @param maxFactWait - maximum time to await the fact
     * @param actual - the actual value at the time of comparison
     * @param expected - the expected value considered earlier
     */
    public void isSame(String assertedFact,
            Object actual, Object expected) {
        condition(assertedFact, ui.condition.isSame(actual, expected));
    }

    /**
     * State Assertion Method that asserts a given fact by checking for the
     * condition:<br><b>
     * is the specified "actual" object NOT SAME as the specified "expected"
     * object
     * </b><br> The comparison is done using the method
     * {@link CommonUtils#smartEqual(Object, Object)}, which makes it very
     * flexible and comfortable for use.<br>
     * The condition is created using {@link SubToolCondition
     * #notSame(Object, Object)} and is awaited to for "maxFactWait" ms.<br>
     * The output of the assertion is handled through the auditing tool:
     * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
     * @param assertedFact - text describing the fact to be asserted
     * @param maxFactWait - maximum time to await the fact
     * @param actual - the actual value at the time of comparison
     * @param expected - the expected value considered earlier
     */
    public void notSame(String assertedFact,
            Object actual, Object expected) {
        condition(assertedFact, ui.condition.notSame(actual, expected));
    }

    /**
     * State Assertion Method that asserts a given fact by checking for the
     * condition:<br><b>
     * IS the specified "value" object TRUE
     * </b><br> The comparison is done using the method
     * {@link CommonUtils#smartEqual(Object, Object)}, which makes it very
     * flexible and comfortable for use.<br>
     * The condition is created using {@link SubToolCondition
     * #isTrue(Object)} and is awaited to for "maxFactWait" ms.<br>
     * The output of the assertion is handled through the auditing tool:
     * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
     * @param assertedFact - text describing the fact to be asserted
     * @param maxFactWait - maximum time to await the fact
     * @param value - the value to be evaluated
     */
    public void isTrue(String assertedFact, Object value) {
        condition(assertedFact, ui.condition.isTrue(value));
    }

    /**
     * State Assertion Method that asserts a given fact by checking for the
     * condition:<br><b>
     * is the specified "value" object NOT TRUE
     * </b><br> The comparison is done using the method
     * {@link CommonUtils#smartEqual(Object, Object)}, which makes it very
     * flexible and comfortable for use.<br>
     * The condition is created using {@link SubToolCondition
     * #notTrue(Object)} and is awaited to for "maxFactWait" ms.<br>
     * The output of the assertion is handled through the auditing tool:
     * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
     * @param assertedFact - text describing the fact to be asserted
     * @param maxFactWait - maximum time to await the fact
     * @param value - the value to be evaluated
     */
    public void notTrue(String assertedFact, Object value) {
        condition(assertedFact, ui.condition.notTrue(value));
    }

    /**
     * State Assertion Method that asserts a given fact by checking for the
     * condition given as input parameter<br>
     * The condition is awaited to for "maxFactWait" ms.<br>
     * The output of the assertion is handled through the auditing tool:
     * {@link SubToolAudit#assertionOutcome(Throwable, String, String)}
     * @param assertedFact - text describing the fact to be asserted
     * @param maxFactWait - maximum time to await the fact
     * @param condition - the used to check the asserted fact
     */
    public void condition(String assertedFact, Condition condition) {
        FailureCauseHolder fch = new FailureCauseHolder();
        try {
            condition.lazyVerify();
        } catch (Exception e) {
            fch.setCause(e);
        }
        ui.audit.assertionOutcome(
                fch.getCause(), assertedFact, condition.toString());
    }
}