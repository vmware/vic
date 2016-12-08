package com.vmware.suitaf.apl;

/**
 * This enumeration represents a classification of the screen component
 * properties. It is intended for used in the following cases:<br>
 * <li> For component logging filters in calls to method {@link
 * AutomationPlatformLinkExt#dumpComponents(ComponentID, int,
 * ComponentDumpLevels, Category...)}.
 * <li> For grouping of items in the enumeration {@link Property}
 * <br><br>
 * @author dkozhuharov
 */
public enum Category {
    ID,
    POSITIONING,
    ACCESSIBILITY,
    DRAWING,
    DRAWINGCOLORS,
    STYLE,
    MOUSE,
    TEXTFORMATTING,
    INFO,
    TABULAR,
    SELECTION,
    AUTOMATION,
    ANY,
    BRIEF,
    REGULAR,
    SYSTEM
}