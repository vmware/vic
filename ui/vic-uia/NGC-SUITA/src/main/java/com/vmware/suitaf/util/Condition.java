package com.vmware.suitaf.util;


/**
 * Abstract class containing the base logic for postponed assert execution.
 * This class introduces a new design pattern. It became necessity with the
 * increase of complexity of the assertion logic.<br>
 * An {@link Condition} implementation will hold the logic for the assertion
 * on a UI condition.
 *
 * @author dkozhuharov
 */
public abstract class Condition {
    private final String description;
    private Boolean lastOutcome = null;

    public Condition(String checkDescription) {
        this.description = checkDescription;
    }

    /**
     * This method awaits the the assertion fact to become <b>true</b> for
     * given amount of milliseconds and returns <b>true</b> if it does.
     * @return <b>true</b> if the assertion fact verifies.
     */
    public boolean await(long maxWait) {
        long start = System.currentTimeMillis();
        long elapsed = 0;

        // This loop must run at least once even if waitTimeout == 0
        while (elapsed <= maxWait) {
            // Check if the condition's criterion is met
            if ( this.estimate() ) {
                break;
            }
            // Wait awhile to allow the expected event to happen
            // The "min" function allow waiting times less that 1000 ms.
            int sleep = Math.min( (int)(maxWait - elapsed), 1000 );
            CommonUtils.sleep( sleep );

            // Recalculate the time it took so far
            elapsed = System.currentTimeMillis() - start;
        }

        return lastOutcome;
    }

    /**
     * This method estimates the assertion fact and returns <b>true</b> if
     * it verifies.
     * @return <b>true</b> if the assertion fact verifies.
     */
    public boolean estimate() {
        lastOutcome = checkImpl();

        return lastOutcome;
    }

    /**
     * This method estimates the assertion fact and throws an exception if
     * it does not verify.
     * @throws AssertionFail if the assertion fact does not verify.
     */
    public void verify() {
        estimate();
        if (!lastOutcome) {
            throw new AssertionFail(toString());
        }
    }

    public void lazyVerify() {
        if (isUnchecked()) {
            estimate();
        }
        if (!lastOutcome) {
            throw new AssertionFail(toString());
        }
    }

    public boolean isTrue() {
        return (lastOutcome == true);
    }

    public boolean isFalse() {
        return (lastOutcome == false);
    }

    public boolean isUnchecked() {
        return (lastOutcome == null);
    }

    @Override
    public String toString() {
        return (lastOutcome==null? "*": (lastOutcome? "+": "-")) + description;
    }

    /**
     * This method holds the implementation for the condition check.
     * @return <b>true</b> if the condition was found, <b>false</b> otherwise
     */
    protected abstract boolean checkImpl();
}