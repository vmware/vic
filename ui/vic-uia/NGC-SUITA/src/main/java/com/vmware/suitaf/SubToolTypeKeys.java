package com.vmware.suitaf;

import static com.vmware.suitaf.SUITA.Factory.apl;
import static com.vmware.suitaf.apl.IDGroup.toIDGroup;

import java.util.Arrays;

import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.ComponentMatcher;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;
import com.vmware.suitaf.apl.Key;

/**
 * This is the sub-tool {@link SubToolTypeKeys} with the following specs:
 * <li> <b>Function Type:</b> ACTIONS simulated at the host OS
 * <li> <b>Description:</b> Typing regular and special keys through the
 *  events of the OS.
 * <li> <b>Based on APL:</b>
 * {@link AutomationPlatformLink#typeKeys(ComponentMatcher, int, Object...)}
 * <li> <b>Auxiliary SubTools:</b>
 * {@link SubToolMouse#down(Object, IDPart...)},
 * {@link SubToolMouse#up()},
 * {@link SubToolAudit#aplFailure(Throwable)}
 */
public class SubToolTypeKeys extends BaseSubTool {
    public SubToolTypeKeys(UIAutomationTool uiAutomationTool) {
        super(uiAutomationTool);
    }

    /**
     * Action Method that simulates pressing of the {@link Key#TAB}
     * key in the host OS.
     */
    public void tab() {
        typeKeysWrapper(null, 0, Key.TAB);
    }

    /**
     * Action Method that simulates pressing of the {@link Key#ENTER}
     * key in the host OS.
     */
    public void enter() {
        typeKeysWrapper(null, 0, Key.ENTER);
    }

    /**
     * Action Method that simulates pressing of the {@link Key#DOWN}
     * key in the host OS.
     */
    public void down() {
        typeKeysWrapper(null, 0, Key.DOWN);
    }

    /**
     * Action Method that simulates pressing of the {@link Key#UP}
     * key in the host OS.
     */
    public void up() {
        typeKeysWrapper(null, 0, Key.UP);
    }

    /**
     * Action Method that simulates setting the focus to the text-entry
     * component, by clicking it and then simulates text-entry through
     * keyboard typing. Both actions are simulated at the level of the host OS.
     * @param text - a string with type-able symbols
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     */
    public void text(String text, Object id, IDPart...propertyFilters) {
        typeKeysWrapper(toIDGroup(id, propertyFilters), 10, text);
    }

    /**
     * Helper Method that wraps the calls to the {@link AutomationPlatformLink
     * #typeKeys(ComponentMatcher, int, Object...)} interface method. It
     * provides error handling and is otherwise transparent.
     * @param id (Optional) - {@link IDGroup} type identifier of the component
     * that will be set under focus before key-typing. The focus setting
     * will be simulated at the FLASH PLAYER level.
     * @param typeDelay - delay in milliseconds between single key types
     * @param keySequence - sequence of {@link Key} constants and strings of
     * type-able characters.
     */
    private void typeKeysWrapper(IDGroup id,
            int typeDelay, Object...keySequence) {
        try {
            ComponentMatcher cm = apl().getComponentMatcher(id);
            apl().typeKeys(cm, typeDelay, keySequence);
        } catch (Exception e) {
            _logger.error("Failed type of " +
                    Arrays.asList(keySequence).toString());
            ui.audit.aplFailure(e);
        }
    }

    /**
     * This utility method calls UI.typeKyes.down() a specified number of times.
     * @param numberOfDowns - number down() to be called.
     */
    public void down(int numberOfDowns){
        for(int i=1; i<=numberOfDowns; i++){
            down();
        }
    }
}
