/**
 *
 */
package com.vmware.suitaf;

import static com.vmware.suitaf.apl.IDGroup.toIDGroup;
import static com.vmware.suitaf.SUITA.Factory.apl;

import java.util.List;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.ComponentMatcher;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;
import com.vmware.suitaf.apl.MouseComponentAction;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.util.Condition;
import com.vmware.suitaf.util.CommonUtils;

/**
 * This is the sub-tool {@link SubToolComponent} with the following specs: <li>
 * <b>Function Type:</b> ACTIONS simulated at the FLASH PLAYER, STATE RETRIEVAL
 * <li><b>Description:</b> All component related functions like get-properties,
 * set-value and click. <li><b>Based on APL:</b>
 * {@link AutomationPlatformLink#getExistingCount(ComponentMatcher)}
 * {@link AutomationPlatformLink#getSingleProperty(ComponentMatcher, Property)}
 * {@link AutomationPlatformLink#getGridProperty(ComponentMatcher, Property)}
 * {@link AutomationPlatformLink #mouseOnComponent(ComponentMatcher, MouseComponentAction, Condition)}
 * {@link AutomationPlatformLink#setValue(ComponentMatcher, String, Condition)}
 * {@link AutomationPlatformLink #setValueByIndex(ComponentMatcher, int, Condition)}
 * <li><b>Auxiliary APL:</b>
 * {@link AutomationPlatformLink#getComponentMatcher(com.vmware.suitaf.apl.IDGroup...)}
 * <li><b>Auxiliary SubTools:</b>
 * {@link SubToolCondition#isFound(Object, IDPart...)},
 * {@link SubToolAudit#aplFailure(Throwable)}
 */
public class SubToolComponent extends BaseSubTool {

    /**
     * This is the sub-tool {@link ComponentProperty} with the following specs:
     * <li> <b>Function Type:</b> Retrieval of property values for a given UI component
     * <li> <b>Description:</b> A set of different get methods to obtain the value
     * of the needed property in whatever type we need
     * <li> <b>Based on:</b>
     * {@link ComponentProperty}
     */
    public final ComponentProperty property;

    /**
     * This is the sub-tool {@link PropertyValue} with the following specs:
     * <li> <b>Function Type:</b> Retrieval and setting of the UI component value
     * <li> <b>Description:</b> Retrieval and setting of the UI component value.
     * for all components that have value properties defined
     * <li> <b>Based on:</b>
     * {@link PropertyValue}
     */
    public final PropertyValue value;

    public SubToolComponent(UIAutomationTool uiAutomationTool) {
        super(uiAutomationTool);;
        this.property = new ComponentProperty();
        this.value = new PropertyValue();
    }


    public class ComponentProperty {
        /**
         * State Retrieval Method that extracts the value of a given property for a
         * UI component.
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested casted to the type of object
         *        requested or <b>null</b> if not supported
         */
        public String get(Property property, Object id,
                IDPart... propertyFilters) {
            ui.condition.isFound(id, propertyFilters).await(
                    SUITA.Environment.getPageLoadTimeout());

            String value = null;
            try {
                ComponentMatcher cm =
                        apl().getComponentMatcher(toIDGroup(id, propertyFilters));
                value = apl().getSingleProperty(cm, property);
            } catch (Exception e) {
                ui.audit.aplFailure(e);
            }
            return value;
        }

        /**
         * Method used to obtain the value of a given property for a
         * given UI component. The returned type depends on the class type
         * provided as first parameter of this method
         *
         * @param type
         *            - the class type of the object that needs to be returned
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested or
         *         exception will be thrown if conversion is not supported
         */
        public <T> T getTyped(Class<T> type, Property property,
                Object id, IDPart... propertyFilters) {
            T value = null;
            try {
                value = CommonUtils.valueOf(
                        type, get(property, id, propertyFilters)
                );

            } catch (Exception e) {
                ui.audit.aplFailure(e);
            }
            return value;
        }

        /**
         * Returns the current value of a property for the UI component as a Boolean object .
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as boolean
         */
        public Boolean getBoolean(Property property, Object id,
                IDPart... propertyFilters) {
            return getTyped (Boolean.class, property, id, propertyFilters);
        }

        /**
         * Returns the current value of a property for the UI component as a Integer object .
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Integer
         */
        public Integer getInteger(Property property, Object id,
                IDPart... propertyFilters) {

            return getTyped (Integer.class, property, id, propertyFilters);
        }

        /**
         * Returns the current value of a property for the UI component as a Long object .
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Long
         */
        public Long getLong(Property property, Object id,
                IDPart... propertyFilters) {
            return getTyped (Long.class, property, id, propertyFilters);
        }


        /**
         * Returns the current value of a property for the UI component as a Double object .
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Double
         */
        public Double getDouble(Property property, Object id,
                IDPart... propertyFilters) {
            return getTyped (Double.class, property, id, propertyFilters);
        }

        /**
         * Returns the current value of a property for the UI component as a Float object .
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Float
         */
        public Float getFloat(Property property, Object id,
                IDPart... propertyFilters) {
            return getTyped (Float.class, property, id, propertyFilters);
        }

        /**
         * Returns the current value of a property for the UI component as a Short object .
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Short
         */
        public Short getShort(Property property, Object id,
                IDPart... propertyFilters) {
            return getTyped (Short.class, property, id, propertyFilters);
        }

        /**
         * Returns the current value of a property for the UI component as a Byte object .
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Byte
         */
        public Byte getByte(Property property, Object id,
                IDPart... propertyFilters) {
            return getTyped (Byte.class, property, id, propertyFilters);
        }

    }

    public class PropertyValue {


        /**
         * Returns the value of the UI component as String object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as String
         */
        public String get(Object id, IDPart... propertyFilters) {

            String value = null;
            try {
                value = ui.component.property.get(Property.VALUE, id, propertyFilters);
            } catch (Exception e) {
                ui.audit.aplFailure(e);
            }
            return value;
        }

        /**
         * Action Method that sets the value of a value holding component.
         *
         * @param valueProperty
         *            - DataProperty object holding the value to be set to the
         *            component. When the DataProperty object holds a
         *            user-defined type of value, the user-defined type needs to
         *            have predefined toString() method that returns the string
         *            representation of the actual value to be set to the UI
         *            component.
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id
         *      values.
         */
        public void set(DataProperty<?> valueProperty, Object id,
                IDPart... propertyFilters) {
            set(valueProperty.stringValue(), id, propertyFilters);
        }

        /**
         * Action Method that sets the value of a value holding component <br>
         * <br>
         *
         * @param value
         *            - value to be set on the component. The framework will call toString()
         *              method of the object as Selenium Flex API works through string serialization.
         *              Be sure to implement proper toString() methods for complex objects.
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         */
        public void set(Object value, Object id, IDPart... propertyFilters) {
            ui.condition.isFound(id, propertyFilters).await(
                    SUITA.Environment.getPageLoadTimeout());

            try {
                ComponentMatcher cm =
                        apl().getComponentMatcher(toIDGroup(id, propertyFilters));
                apl().setValue(cm, value.toString(), null);
            } catch (Exception e) {
                ui.audit.aplFailure(e);
            }
        }


        /**
         * Method used to obtain the value of for a given UI component.
         * The returned type depends on the class type provided as
         * first parameter of this method
         *
         * @param type
         *            - the class type of the object that needs to be returned
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the component or
         *         exception will be thrown if conversion is not supported
         */
        public <T> T getTyped(Class<T> type, Object id, IDPart... propertyFilters) {
            T value = null;
            try {
                value = CommonUtils.valueOf(
                        type, get(id, propertyFilters)
                );

            } catch (Exception e) {
                ui.audit.aplFailure(e);
            }
            return value;
        }

        /**
         * Returns the value of the UI Component as Boolean Object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Boolean
         */
        public Boolean getBoolean(Object id, IDPart... propertyFilters) {
            return getTyped(Boolean.class, id, propertyFilters);
        }

        /**
         * Returns the value of the UI Component as Integer Object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Integer
         */
        public Integer getInteger(Object id, IDPart... propertyFilters) {
            return getTyped(Integer.class, id, propertyFilters);
        }

        /**
         * Returns the value of the UI Component as Long Object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Long
         */
        public Long getLong(Object id, IDPart... propertyFilters) {
            return getTyped(Long.class, id, propertyFilters);
        }

        /**
         * Returns the value of the UI Component as Short Object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Short
         */
        public Short getShort(Object id, IDPart... propertyFilters) {
            return getTyped(Short.class, id, propertyFilters);
        }

        /**
         * Returns the value of the UI Component as Byte Object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Byte
         */
        public Byte getByte(Object id, IDPart... propertyFilters) {
            return getTyped(Byte.class, id, propertyFilters);
        }


        /**
         * Returns the value of the UI Component as Float Object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Float
         */
        public Float getFloat(Object id, IDPart... propertyFilters) {
            return getTyped(Float.class, id, propertyFilters);
        }


        /**
         * Returns the value of the UI Component as Double Object
         *
         * @param property
         *            - the property to the retrieved
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return the value of the property requested as Double
         */
        public Double getValueAsDouble(Object id, IDPart... propertyFilters) {
            return getTyped(Double.class, id, propertyFilters);
        }

        /**
         * State Retrieval Method that gets the selected grid row as a string
         * concatenation <br>
         * <br>
         *
         * @param id
         *            - the base id of the component
         * @param propertyFilters
         *            - additional property filters to the base id
         * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
         * @return concatenation of the grid cells "toString()" invocation
         */
        public String getSelected(Object id, IDPart... propertyFilters) {
            ui.condition.isFound(id, propertyFilters).await(
                    SUITA.Environment.getPageLoadTimeout());

            String value = null;
            try {
                ComponentMatcher cm =
                        apl().getComponentMatcher(toIDGroup(id, propertyFilters));
                value = apl().getSingleProperty(cm, Property.VALUE);
            } catch (Exception e) {
                ui.audit.aplFailure(e);
            }
            return value;
        }

    }


    /**
     * State Retrieval Method that checks object existence. <br>
     * <br>
     *
     * @param id
     *            - the base id of the component
     * @param propertyFilters
     *            - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return <b>true</b> if component matching the ID and compliant to all
     *         property filters has been found.
     */
    public boolean exists(Object id, IDPart... propertyFilters) {
        Boolean value = null;
        try {
            ComponentMatcher cm =
                    apl().getComponentMatcher(toIDGroup(id, propertyFilters));
            value = (apl().getExistingCount(cm) > 0);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }

        return value;
    }
    
    public int existingCount(Object id, IDPart... propertyFilters) {
        try {
            ComponentMatcher cm =
                    apl().getComponentMatcher(toIDGroup(id, propertyFilters));
            return apl().getExistingCount(cm);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }
        return 0;
    }
    
    /**
     * Action Method that clicks on an object <br>
     * <br>
     *
     * @param id
     *            - the base id of the component
     * @param propertyFilters
     *            - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     */
    public void click(Object id, IDPart... propertyFilters) {
        click(null, id, propertyFilters);
    }

    /**
     * Action Method that clicks an object using {@link Condition} for
     * validation <br>
     * <br>
     *
     * @param asserter
     *            - {@link Condition} for validation
     * @param id
     *            - the base id of the component
     * @param propertyFilters
     *            - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     */
    public void click(Condition asserter, Object id, IDPart... propertyFilters) {
        ui.condition.isFound(id, propertyFilters).await(
                SUITA.Environment.getPageLoadTimeout());

        try {
            ComponentMatcher cm =
                    apl().getComponentMatcher(toIDGroup(id, propertyFilters));
            apl().mouseOnComponent(cm, MouseComponentAction.CLICK, asserter);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }
    }


    /**
     * Action Method that sets the focus to the focusable component and makes it
     * accept any further key types <br>
     * <br>
     *
     * @param id
     *            - the base id of the component
     * @param propertyFilters
     *            - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     */
    public void setFocus(Object id, IDPart... propertyFilters) {
        ui.condition.isFound(id, propertyFilters).await(
                SUITA.Environment.getPageLoadTimeout());

        try {
            ComponentMatcher cm =
                    apl().getComponentMatcher(toIDGroup(id, propertyFilters));
            apl().setFocus(cm);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }
    }

    /**
     * State Retrieval Method that finds the row index of an item in a grid <br>
     * <br>
     *
     * @param item
     *            - the item to be reached
     * @param id
     *            - the base id of the component
     * @param propertyFilters
     *            - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     * @return - the row number where the item was first found or -1 otherwise.
     */
    public int indexOfItemInGrid(String item, Object id,
            IDPart... propertyFilters) {
        ui.condition.isFound(id, propertyFilters).await(
                SUITA.Environment.getPageLoadTimeout());

        List<List<String>> gridData = null;
        try {
            ComponentMatcher cm =
                    apl().getComponentMatcher(toIDGroup(id, propertyFilters));
            gridData = apl().getGridProperty(cm, Property.GRID_DATA);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }

        for (int i = 0; i < gridData.size(); i++) {
            if (gridData.get(i).contains(item)) {
                return i;
            }
        }
        return -1;
    }

    /**
     * Action Method that selects a row in the grid by index <br>
     * <br>
     *
     * @param index
     *            - row number to be selected
     * @param id
     *            - the base id of the component
     * @param propertyFilters
     *            - additional property filters to the base id
     * @see {@link IDGroup#toIDGroup(Object, IDPart...)} for valid id values.
     */
    public void selectByIndex(int index, Object id, IDPart... propertyFilters) {
        ui.condition.isFound(id, propertyFilters).await(
                SUITA.Environment.getPageLoadTimeout());

        try {
            ComponentMatcher cm =
                    apl().getComponentMatcher(toIDGroup(id, propertyFilters));
            apl().setValueByIndex(cm, index, null);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }
    }




}