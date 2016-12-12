package com.vmware.suitaf.apl;

/**
 * These dump-levels definitions determines for which components a detailed
 * property info will be logged.
 *
 * @author dkozhuharov
 */
public enum ComponentDumpLevels {
    NO_DETAILS("No components are logged in detail"),
    ROOT_ONLY("Only the starting component is logged in detail"),
    USER_INPUT_ONLY("Only user input components are logged in detail"),
    ALL_DETAILS("All components are logged in detail");

    public final String description;
    ComponentDumpLevels(String description) {
        this.description = description;
    }
}
