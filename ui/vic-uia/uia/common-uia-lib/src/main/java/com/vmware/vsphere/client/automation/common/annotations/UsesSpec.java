package com.vmware.vsphere.client.automation.common.annotations;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import com.vmware.vim.binding.vmodl.link;

/**
 * Annotation for a fields in class associated with a Spec classes and have to
 * be assigned from the WorkflowSpec
 *
 * The initialization will happen automatically by a {@link link
 * WorkflowSpecInitializer} call in the corresponding appropriate phase of the
 * life cycle of the class in which the annotation is present
 *
 * <h1>
 * Sample usage:</h1>
 *
 * <pre>
 * &#064;WorkflowSpecAnnotation
 * private SomeSpecClass specVariableName;
 * </pre>
 *
 * <pre>
 * &#064;WorkflowSpecAnnotation
 * private SomeSpecClass[] arrayOfSpecsVariableName;
 * </pre>
 *
 * <pre>
 * &#064;WorkflowSpecAnnotation
 * private List&lt;SomeSpecClass&gt; listOfSpecsVariableName;
 * </pre>
 */
@Retention(RetentionPolicy.RUNTIME)
@Target(value = ElementType.FIELD)
public @interface UsesSpec {

}
