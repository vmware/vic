/* **********************************************************************
 * Copyright 2010 VMware, Inc.  All rights reserved. VMware Confidential
 * *********************************************************************/

package com.vmware.suitaf;

import static com.vmware.suitaf.SUITA.Factory.apl;
import static com.vmware.suitaf.apl.IDGroup.toIDGroup;

import java.awt.Point;

import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.ComponentMatcher;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;
import com.vmware.suitaf.apl.MouseButtonAction;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.util.Condition;

/**
 * This is the sub-tool {@link SubToolScreenMouse} with the following specs:
 * <li> <b>Function Type:</b> ACTIONS simulated at the host OS
 * <li> <b>Description:</b> Execution of mouse commands using the coordinate
 * system of screen through the events of the OS.
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
public class SubToolScreenMouse extends BaseSubTool {
    private static final Long MOUSE_ACT_COOLDOWN = 500L;
    private static final Integer CATCH_POINT_OFFSET = 10;

    public SubToolScreenMouse(UIAutomationTool uiAutomationTool) {
        super(uiAutomationTool);
    }

    /**
     *
     * State Retrieval Method that gets the "catch point" of a UI component.
     * The "catch point" is 10 pixels below the middle of the upper edge of
     * the component. This point is suitable to be used for clicking or
     * dragging of that component.
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return a {@link Point} object representing the "catch point"
     */
    public Point getCatchPoint(Object id, IDPart...propertyFilters) {

        ui.condition.isFound(id, propertyFilters).await(
                SUITA.Environment.getPageLoadTimeout());

        Point catchPoint = null;
        try {
            ComponentMatcher cm = apl().getComponentMatcher(
                    toIDGroup(id, propertyFilters));
            int x = Integer.parseInt(
                    apl().getSingleProperty(cm, Property.SCR_X));
            int y = Integer.parseInt(
                    apl().getSingleProperty(cm, Property.SCR_Y));
            int w = Integer.parseInt(
                    apl().getSingleProperty(cm, Property.POS_WIDTH));

            catchPoint = new Point(x + w/2, y + CATCH_POINT_OFFSET);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }

        return catchPoint;
    }

    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> press-down of the mouse button on the current mouse pointer position
     */
    public void down() {
        act(null, MouseButtonAction.LEFT_DOWN, MOUSE_ACT_COOLDOWN);
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> moves the mouse pointer to a specified point on screen
     * <li> press-down of the mouse button on the current mouse pointer position
     * @param p - the point for the action; the coordinates must be relative to
     * the Flex application's (its upper-left corner).
     */
    public void down(Point p) {
        act(null, p, MOUSE_ACT_COOLDOWN,
                MouseButtonAction.LEFT_DOWN, MOUSE_ACT_COOLDOWN);
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> retrieves the "catch point" of a specified UI component
     * <li> moves the mouse pointer to the "catch point"
     * <li> press-down of the mouse button on the current mouse pointer position
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return the "catch point" to which the action was applied
     */
    public Point down(Object id, IDPart...propertyFilters) {
        return down((Condition) null, id, propertyFilters);
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> retrieves the "catch point" of a specified UI component
     * <li> moves the mouse pointer to the "catch point"
     * <li> press-down of the mouse button on the current mouse pointer position
     * <li> evaluates the success of the action using specified "success"
     * {@link Condition}
     * <li> if the action was unsuccessful - an attempt is made to retry it.
     * @param asserter (Optional) - a {@link Condition} object that could
     * evaluate if the action was a "success"
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return the "catch point" to which the action was applied
     */
    public Point down(Condition asserter,
            Object id, IDPart...propertyFilters) {
        Point p = getCatchPoint(id, propertyFilters);
        act(asserter, p, MOUSE_ACT_COOLDOWN,
                MouseButtonAction.LEFT_DOWN, MOUSE_ACT_COOLDOWN);
        return p;
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> release-up of the mouse button on the current mouse pointer position
     */
    public void up() {
        act(null, MouseButtonAction.LEFT_UP, MOUSE_ACT_COOLDOWN);
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> moves the mouse pointer to a specified point on screen
     * <li> release-up of the mouse button on the current mouse pointer position
     * @param p - the point for the action; the coordinates must be relative to
     * the Flex application's (its upper-left corner).
     */
    public void up(Point p) {
        act(null, p, MOUSE_ACT_COOLDOWN,
                MouseButtonAction.LEFT_UP, MOUSE_ACT_COOLDOWN);
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> retrieves the "catch point" of a specified UI component
     * <li> moves the mouse pointer to the "catch point"
     * <li> release-up of the mouse button on the current mouse pointer position
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return the "catch point" to which the action was applied
     */
    public Point up(Object id, IDPart...propertyFilters) {
        return up((Condition) null, id, propertyFilters);
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> retrieves the "catch point" of a specified UI component
     * <li> moves the mouse pointer to the "catch point"
     * <li> release-up of the mouse button on the current mouse pointer position
     * <li> evaluates the success of the action using specified "success"
     * {@link Condition}
     * <li> if the action was unsuccessful - an attempt is made to retry it.
     * @param asserter (Optional) - a {@link Condition} object that could
     * evaluate if the action was a "success"
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return the "catch point" to which the action was applied
     */
    public Point up(Condition asserter,
            Object id, IDPart...propertyFilters) {
        Point p = getCatchPoint(id, propertyFilters);
        act(asserter, p, MOUSE_ACT_COOLDOWN,
                MouseButtonAction.LEFT_UP, MOUSE_ACT_COOLDOWN);
        return p;
    }

    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> moves the mouse pointer to a specified point on screen
     * @param p - the point for the action; the coordinates must be relative to
     * the Flex application's (its upper-left corner).
     */
    public void move(Point p) {
        act(null, p, MOUSE_ACT_COOLDOWN);
    }

    /**
     *
     * Action Method that simulates complex mouse actions. It accepts an
     * arbitrary mouse action sequence consisting of:
     * <li> {@link MouseButtonAction} constants - triggers execution of a mouse
     * button action;
     * <li> {@link Point} objects - triggers a mouse-move action;
     * <li> {@link Long} values - initiates waiting for the given number of
     * milliseconds - used for mouse action cool-down delays;
     * <br>
     * After the whole action sequence was executed:
     * <li> evaluates the success of the action sequence using the specified
     * "success" {@link Condition}
     * <li> if the action sequence was unsuccessful - an attempt is made to
     * retry it.
     * @param asserter (Optional) - a {@link Condition} object that could
     * evaluate if the action was a "success"
     * @param mouseActionsSequence - the sequence of mouse action elements
     */
    public void act(Condition asserter, Object... mouseActionsSequence) {
        try {
            apl().mouseOnScreen(asserter, mouseActionsSequence);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }
    }
    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> retrieves the "catch point" of a specified UI component
     * <li> moves the mouse pointer to the "catch point"
     * <li> press-down of the mouse button on the current mouse pointer position
     * <li> release-up of the mouse button on the current mouse pointer position
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return the "catch point" to which the action was applied
     */
    public Point click(Object id, IDPart...propertyFilters) {
        return click((Condition) null, id, propertyFilters);
    }

    /**
     * Action Method that simulates the following mouse actions at the host OS:
     * <li> retrieves the "catch point" of a specified UI component
     * <li> moves the mouse pointer to the "catch point"
     * <li> press-down of the mouse button on the current mouse pointer position
     * <li> evaluates the success of the action using specified "success"
     * {@link Condition}
     * <li> if the action was unsuccessful - an attempt is made to retry it.
     * @param asserter (Optional) - a {@link Condition} object that could
     * evaluate if the action was a "success"
     * @param id - the base id of the component
     * @param propertyFilters - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return the "catch point" to which the action was applied
     */
    public Point click(Condition asserter,
            Object id, IDPart...propertyFilters) {
        Point p = getCatchPoint(id, propertyFilters);
        act(asserter, p, MOUSE_ACT_COOLDOWN,
                MouseButtonAction.LEFT_DOWN, MOUSE_ACT_COOLDOWN,
                MouseButtonAction.LEFT_UP, MOUSE_ACT_COOLDOWN);
        return p;
    }
}
