package com.vmware.suitaf;

import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.AutomationPlatformLinkExt;
import com.vmware.suitaf.apl.Category;
import com.vmware.suitaf.apl.ComponentDumpLevels;
import com.vmware.suitaf.apl.ComponentMatcher;
import com.vmware.suitaf.apl.HostLogger;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;
import com.vmware.suitaf.apl.MouseComponentAction;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.apl.SpecialStateHandler;
import com.vmware.suitaf.apl.SpecialStates;
import com.vmware.suitaf.util.Condition;

/**
 * This class serves as facade for the set of UI automation tools included
 * in the SUITA framework. It provides simple intuitive access to the needed
 * functionality.
 *
 * @author dkozhuharov
 *
 */
public final class UIAutomationTool {
    // =====================================================================
    // Initialization Section of the UI tool
    // =====================================================================
    /**
     * This is a convenience shortcut to the logger of the SUITA framework.
     */
    public final HostLogger logger;
    /**
     * This is the sub-tool {@link SubToolCondition} with the following specs:
     * <li> <b>Function Type:</b> COMPONENT FACTORY
     * <li> <b>Description:</b> Creation of {@link Condition} instances to be used
     * for checking, awaiting or assertion of specific conditions of the UI or
     * of the test application.
     * <li> <b>Based on SubTool:</b>
     * {@link SubToolComponent#exists(Object, IDPart...)}
     */
    public final SubToolCondition condition;
    /**
     * This is the sub-tool {@link SubToolAssertFor} with the following specs:
     * <li> <b>Function Type:</b> STATE RETRIEVAL
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
    public final SubToolAssertFor assertFor;
    /**
     * This is the sub-tool {@link SubToolComponent} with the following specs:
     * <li> <b>Function Type:</b> ACTIONS simulated at the FLASH PLAYER,
     * STATE RETRIEVAL
     * <li> <b>Description:</b> All component related functions like get-properties,
     * set-value and click.
     * <li> <b>Based on APL:</b>
     * {@link AutomationPlatformLink#getExistingCount(ComponentMatcher)}
     * {@link AutomationPlatformLink#getSingleProperty(ComponentMatcher, Property)}
     * {@link AutomationPlatformLink#getGridProperty(ComponentMatcher, Property)}
     * {@link AutomationPlatformLink
     * #mouseOnComponent(ComponentMatcher, MouseComponentAction, Condition)}
     * {@link AutomationPlatformLink#setValue(ComponentMatcher, String, Condition)}
     * {@link AutomationPlatformLink
     * #setValueByIndex(ComponentMatcher, int, Condition)}
     * <li> <b>Auxiliary APL:</b>
     * {@link AutomationPlatformLink#getComponentMatcher(
     * com.vmware.suitaf.apl.IDGroup...)}
     * <li> <b>Auxiliary SubTools:</b>
     * {@link SubToolCondition#isFound(Object, IDPart...)},
     * {@link SubToolAudit#aplFailure(Throwable)}
     */
    public final SubToolComponent component;
    /**
     * This is the sub-tool {@link SubToolMouse} with the following specs:
     * <li> <b>Function Type:</b> ACTIONS simulated at the host OS
     * <li> <b>Description:</b> Execution of mouse commands using the coordinate
     * system of the Flex application through the events of the OS.
     * <li> <b>Based on APL:</b>
     * {@link AutomationPlatformLink#mouseOnApplication(Condition, Object...)}
     * <li> <b>Auxiliary APL:</b>
     * {@link AutomationPlatformLink#getSingleProperty(ComponentMatcher, Property)},
     * {@link AutomationPlatformLink#getComponentMatcher(
     * com.vmware.suitaf.apl.IDGroup...)}
     * <li> <b>Auxiliary SubTools:</b>
     * {@link SubToolCondition#isFound(Object, IDPart...)},
     * {@link SubToolAudit#aplFailure(Throwable)}
     */
    public final SubToolMouse mouse;

    /**
     * This is the sub-tool {@link SubToolScreenMouse} with the following specs:
     * <li> <b>Function Type:</b> ACTIONS simulated at the host OS
     * <li> <b>Description:</b> Execution of mouse commands using the coordinate
     * system of the screen through the events of the OS.
     * <li> <b>Based on APL:</b>
     * {@link AutomationPlatformLink#mouseOnApplication(Condition, Object...)}
     * <li> <b>Auxiliary APL:</b>
     * {@link AutomationPlatformLink#getSingleProperty(ComponentMatcher, Property)},
     * {@link AutomationPlatformLink#getComponentMatcher(
     * com.vmware.suitaf.apl.IDGroup...)}
     * <li> <b>Auxiliary SubTools:</b>
     * {@link SubToolCondition#isFound(Object, IDPart...)},
     * {@link SubToolAudit#aplFailure(Throwable)}
     */
    public final SubToolScreenMouse screenMouse;


    /**
     * This is the sub-tool {@link SubToolBrowser} with the following specs:
     * <li> <b>Function Type:</b> ACTIONS simulated at the host BROWSER
     * <li> <b>Description:</b> Execution of browser commands like open-url,
     * prev-page and next-page using the browser's JavaScript engine.
     * <li> <b>Based on APL:</b>
     * {@link AutomationPlatformLink#openUrl(String)}
     * {@link AutomationPlatformLink#goBack()}
     * {@link AutomationPlatformLink#goForward()}
     * <li> <b>Auxiliary SubTools:</b>
     * {@link SubToolSpecialState#getHandler(SpecialStates)},
     * {@link SubToolAudit#aplFailure(Throwable)}
     */
    public final SubToolBrowser browser;
    /**
     * This is the sub-tool {@link SubToolAudit} with the following specs:
     * <li> <b>Function Type:</b> STATE RETRIEVAL
     * <li> <b>Description:</b> Contains commands that retrieve and log the
     * current state of the UI and of the test-application itself.
     * <li> <b>Based on APL:</b>
     * {@link AutomationPlatformLinkExt#dumpComponents(
     * ComponentMatcher, int, ComponentDumpLevels, Category...)},
     * {@link AutomationPlatformLinkExt#captureBitmap(String)}
     * <li> <b>Auxiliary APL:</b>
     * {@link AutomationPlatformLink#getComponentMatcher(IDGroup...)}
     */
    public final SubToolAudit audit;
    /**
     * This is the sub-tool {@link SubToolSpecialState} with the following specs:
     * <li> <b>Function Type:</b> COMPONENT FACTORY
     * <li> <b>Description:</b> Creation of {@link SpecialStateHandler} instances.
     * <li> <b>Based on APL:</b>
     * {@link AutomationPlatformLink#getSpecialStateHandler(SpecialStates)}
     * <li> <b>Auxiliary SubTools:</b>
     * {@link SubToolAudit#aplFailure(Throwable)}
     */
    public final SubToolSpecialState specialState;
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
    public final SubToolTypeKeys typeKeys;


    UIAutomationTool(HostLogger logger) {
        this.assertFor = new SubToolAssertFor(this);
        this.audit = new SubToolAudit(this);
        this.browser = new SubToolBrowser(this);
        this.component = new SubToolComponent(this);
        this.condition = new SubToolCondition(this);
        this.mouse = new SubToolMouse(this);
        this.screenMouse = new SubToolScreenMouse(this);
        this.logger = logger;
        this.specialState = new SubToolSpecialState(this);
        this.typeKeys = new SubToolTypeKeys(this);
    }
}
