package com.vmware.suitaf;

import static com.vmware.suitaf.SUITA.Factory.apl;

import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.SpecialStateHandler;
import com.vmware.suitaf.apl.SpecialStates;

/**
 * This is the sub-tool {@link SubToolSpecialState} with the following specs:
 * <li> <b>Function Type:</b> COMPONENT FACTORY
 * <li> <b>Description:</b> Creation of {@link SpecialStateHandler} instances.
 * <li> <b>Based on APL:</b>
 * {@link AutomationPlatformLink#getSpecialStateHandler(SpecialStates)}
 * <li> <b>Auxiliary SubTools:</b>
 * {@link SubToolAudit#aplFailure(Throwable)}
 */
public class SubToolSpecialState extends BaseSubTool {
    public SubToolSpecialState(UIAutomationTool uiAutomationTool) {
        super(uiAutomationTool);
    }

    /**
     * Factory Method that creates a handler of special states. I is used to
     * detect and eventually handle special states that could prevent tests from
     * normal execution.
     * @param sState - {@link SpecialStates} designator of the state to be
     * handled.
     * @return an handler instance
     */
    public SpecialStateHandler getHandler(SpecialStates sState) {
        SpecialStateHandler value = null;
        try {
            value = apl().getSpecialStateHandler(sState);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }

        return value;
    }
}
