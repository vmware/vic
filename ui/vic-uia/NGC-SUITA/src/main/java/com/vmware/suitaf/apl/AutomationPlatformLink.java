package com.vmware.suitaf.apl;

import java.awt.Point;
import java.util.List;
import java.util.Map;

import com.vmware.suitaf.util.Condition;

/**
 * This interface is the base for decoupling of the SUITA framework and the
 * platform used for UI test automation (e.g. Silk, Selenium, Sikuli, etc.).
 * The functionally included in the interface covers the needs of the UI test
 * implementation.
 *
 * @author dkozhuharov
 */
public interface AutomationPlatformLink {

    // ==============================================================
    //    (APL-10) Link setup and recovery;
    // ==============================================================

    /**
     * This method is a complex setup and recovery tool. After its execution
     * the following state is expected:
     * <li> A link is established with a running automation agent
     * <li> A browser is started and control is gained on it.
     * <br><br>
     * The configuration data is provided as opened Map of key-value pairs.
     * The information they provide could include the automation agent
     * connection data and also the browser type.<br>
     * To reach the expected state the method could execute one of several
     * scenarios. The particular scenario is chosen automatically based on the
     * current state of the implementor instance and the provided parameters.
     * Here are some of the possibilities:
     * <br><br>
     * <li> Initialization - when no connection with an Automation Agent has
     * been established yet.
     * <li> Switch of Agent - when a current link to an agent exists and
     * provided parameters has different IP/port of the Agent.
     * <li> Switch of Browser - when a current link to an agent exists and
     * provided parameters has same IP/port of the Agent, but different browser.
     * <li> Restart of Browser - when a current link to an agent exists and
     * provided parameters for Agent IP/port and Browser are the same.
     * <br><br>
     * In general these scenarios are only for illustration purposes and are
     * not enforced by the interface. This is deliberately because future
     * Automation Platforms could adopt different processes of connection and
     * control of the browser.
     * <br><br>
     * If the destination state could not be reached (i.e. to have control
     * of a running browser) this method could call special cleanup & recovery
     * procedures. They could include forced cleanup of stalled browsers or
     * even restart of the Automation Agent itself.
     * <br><br>
     * This method is expected to be called at the start of tests or test groups
     * <br><br>
     * @param linkInitParams - a map of key-value pairs that represent a
     * complete and meaningful set of parameter for the initialization of
     * the current APL implementation.
     * @param hostLogger - a call-back interface to allow controlled logging
     * in the main SUITA log stream.
     */
    void resetLink(Map<String, String> linkInitParams, HostLogger hostLogger);

    // ==============================================================
    //    (APL-11) Browser control;
    // ==============================================================

    // Browser page control functions
    /**
     * Opens the specified URL in the browser.
     * @param url - address of the page to be opened by the browser
     */
    void openUrl(String url);
    /**
     * Returns the URL of the currently opened page
     * @return
     */
    String getUrl();
    /**
     * Open next page in the page history list
     */
    void goForward();
    /**
     * Open previous page in the page history list
     */
    void goBack();
    /**
     * Reloads the current page opened in the browser
     */
    void reloadPage();

    // Browser window control functions
    /** This method should implement switching of browser window modes as
     * described by the enumeration {@link WindowModes}
     * <br><br>
     * @param mode - target window mode
     */
    void setBrowserWindowMode(WindowModes mode);
    /** This method should return the current browser window mode as
     * described by the enumeration {@link WindowModes}
     * <br><br>
     * @return
     */
    WindowModes getBrowserWindowMode();

    /**
     * Delivers a native {@link SpecialStateHandler} implementation,
     * @param state - id of the state that needs to be handled
     * @return {@link SpecialStateHandler} instance
     * @throws UnsupportedOperationException - if no appropriate implementation
     *  is available
     */
    SpecialStateHandler getSpecialStateHandler(SpecialStates state);

    // ==============================================================
    //    (APL-13) Screen inspection (parent/child component relations
    //             and property enumeration)
    // ==============================================================

    /**
     * This is a factory method that construct a native {@link ComponentMatcher}
     * implementations. They are built from an array of {@link IDGroup}s. The
     * implementation should ignore any <b>null</b> groups.
     * <br><br>
     * The implementation of the {@code getComponentMatcher} factory method
     * could perform any translation and normalization actions on the general
     * format data to achieve consistent native representation of the component
     * identifier.
     * <br><br>
     * @param idGroups - vararg array of {@link IDGroup} elements. Each instance
     * contains one Property-Value pair for the {@link Property#GROUPROLE} that
     * determines the role of the corresponding entity.
     * @return Native representation of the generic component identifier. This
     * representation is the only component ID format that the methods in this
     * interface accept.
     * @throws NonUsableComponentIDException - fail due to not supported ID part
     * type or non-compliant ID data.
     */
    ComponentMatcher getComponentMatcher(IDGroup...idGroups);

    // ==============================================================
    //    (APL-20) Existence check (single component existence/absence)
    // ==============================================================

    /**
     * This method is the base for methods awaiting appearance/disappearance
     * of given screen component. It is also a must-have before any component
     * operation is launched;
     * <br><br>
     * @param componentID - component identifier to be used to search for
     * components on screen
     * @return - the count of objects existing on screen that match the
     * provided component ID
     */
    int getExistingCount(ComponentMatcher componentID);

    // ==============================================================
    //    (APL-21) Property get
    // ==============================================================

    /**
     * This method's implementation must "understand" and match any of the
     * standardized property names from {@link Property} to the
     * Automation Platform specific properties. The usage of standardized
     * property names allows cross platform mobility of the test
     * implementations.<br>
     * This method works for single- or multiple- valued properties.
     *
     * @param baseComponent - component whose property values will be extracted
     * @param property      - property identifier
     * @return the String list equivalent of the property value(s)
     * @throws RuntimeException in case of any failure extracting the values
     */
    String getSingleProperty(
            ComponentMatcher baseComponent, Property property);

    List<String> getListProperty(
            ComponentMatcher baseComponent, Property property);

    List<List<String>> getGridProperty(
            ComponentMatcher baseComponent, Property property);

    // ==============================================================
    //    (APL-22) State set (user controllable props.)
    // ==============================================================

    /**
     * This method sets the value of a value-holding component.<br>
     * Implementation must be provided for components with fixed list of states
     * and for such with free values (string, integer, etc.). Implementation
     * must support execution for the following generic types of components:
     * <li> text-box         - the text value
     * <li> numeric stepper  - the numeric value
     * <li> date-time picker - the date-time value
     * <li> combo-box - the selected combo-box item
     * <li> check-box - the true(checked)/false(unchecked) state value
     * <li> list      - the selected list item
     * <li> radio-button - the true(selected)/false(unselected) state value
     * <li> tab-control  - the selected tab-name value
     * <li> tree-control - the selected tree-item value
     * <br><br>
     * @param valueholdingComponent - the component whose value will be set
     * @param newValue              - the value to be set
     */
    void setValue( ComponentMatcher valueholdingComponent,
            String newValue, Condition asserter);

    /**
     * This method sets the value of a component by setting its selection index.
     * Implementation must be provide for all components that could have a fixed
     * list of states. Implementation must support execution for the following
     * generic types of components:
     * <li> combo-box
     * <li> grid
     * <li> list
     * <li> radio-button group
     * <li> tab-control
     * <li> tree-control
     *
     * @param indexableComponent - the component whose selection will be changed
     * @param newIndex           - the index of the value to be selected
     */
    void setValueByIndex( ComponentMatcher indexableComponent,
            int newIndex, Condition asserter);


    // ==============================================================
    //    (APL-23) Action invocation (mouse/keys user actions)
    // ==============================================================

    // ==== Keyboard action control methods
    /**
     * Sets the keyboard input focus to specified component. Method should
     * throw exception if the component does not support the operation.
     *
     * @param focusableComponent
     */
    void setFocus(ComponentMatcher focusableComponent);
    /**
     * Types a sequence of keys with optional delay between each key. The key
     * typing goes to the currently focused component. Optionally a focus
     * component could be given. This variant is provided for implementors
     * that could not provide reliable focus shifting behavior.
     * <br><br>
     * Keys to be typed are provided as an array of {@link Object}s. Each item
     * in this array should be either a {@link String} of typeable symbols or
     * {@link Key} enum values for the special keys.
     *
     * @param focusableComponent - (Optional)
     * @param typeDelay
     * @param keySequence
     */
    void typeKeys(ComponentMatcher focusableComponent,
            int typeDelay, Object...keySequence);

    // ==== Mouse action control methods
    /**
     * Method that allows execution of sequence of mouse operation in the
     * context of the screen coordinate system. The sequence could be composed
     * of the following types of elements:
     * <li> {@link Long} - numeric value will be implemented as "sleep" command
     * <li> {@link Point} - couple of coordinates will be implemented as
     * "mouse-move" command. The coordinates will be interpreted as screen
     * related.
     * <li> {@link MouseButtonAction} - mouse action enum values will be implemented
     * as the corresponding mouse button actions.
     *
     * @param asserter - an optional {@link Condition} instance. If provided
     * is used in workaroud loops. Workaroud loops contain several alternative
     * implementations of the action and need the asserter to find when the
     * action is successful.
     * @param mouseActionsSequence - array of mouse action elements. The order
     * of execution will be the parameter order.
     */
    void mouseOnScreen(Condition asserter, Object...mouseActionsSequence);

    /**
     * Method that allows execution of sequence of mouse operation in the
     * context of the application's coordinate system. The sequence could be
     * composed of the following types of elements:
     * <li> {@link Long} - numeric value will be implemented as "sleep" command
     * <li> {@link Point} - couple of coordinates will be implemented as
     * "mouse-move" command. The coordinates will be interpreted as application
     * related. This mouse operation must be consistent with the values of
     * component's positioning properties. Positioning properties for the UI
     * component are: {@link Property#POS_X}, {@link Property#POS_Y},
     * {@link Property#POS_WIDTH}, {@link Property#POS_HEIGHT}
     * <li> {@link MouseButtonAction} - mouse action enum values will be implemented
     * as the corresponding mouse button actions.
     *
     * @param asserter - an optional {@link Condition} instance. If provided
     * is used in workaroud loops. Workaroud loops contain several alternative
     * implementations of the action and need the asserter to find when the
     * action is successful.
     * @param mouseActionsSequence - array of mouse action elements. The order
     * of execution will be the parameter order.
     */
    void mouseOnApplication(Condition asserter, Object...mouseActionsSequence);

    /**
     * This method executes mouse action on a given screen component.
     * The component must be clickable and enabled. The available mouse actions
     * are defined in {@link MouseComponentAction}.
     * <br><br>
     * @param clickableComponent - {@link ComponentMatcher} for the screen
     * component to be clicked;
     * @param action - {@link MouseComponentAction} identifier of the button
     * action to use on the component;
     * @param asserter - an optional {@link Condition} instance. If provided
     * is used in workaroud loops. Workaroud loops contain several alternative
     * implementations of the action and need the asserter to find when the
     * action is successful.
     */
    void mouseOnComponent(ComponentMatcher clickableComponent,
            MouseComponentAction action, Condition asserter);

}
