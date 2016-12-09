package com.vmware.suitaf;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * This is the base class for all sub-tool classes. It only provides the
 * back reference to the common UI tool container class. It is used in each
 * sub-tool to refer and use the functionality of the others.
 * @author dkozhuharov
 */
public class BaseSubTool {

    protected static final Logger _logger = LoggerFactory.getLogger(BaseSubTool.class);

    /**
     * This field contains the back-reference to the UI tool container class.
     */
    protected final UIAutomationTool ui;

    /**
     * This constructor initializes the back reference field
     * @param uiAutomationTool
     */
    public BaseSubTool(UIAutomationTool uiAutomationTool) {
        ui = uiAutomationTool;
    }

    /**
     * This private definition hides the default constructor from all child
     * classes.
     */
    @SuppressWarnings("unused")
    private BaseSubTool() {
        ui = null;
    }
}
