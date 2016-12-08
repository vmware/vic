/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.workflow;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * An annotation used for associating a View class to a Step class.
 * The view's "validate" method is called before the "execute" method of the step.
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface View {
   // The View class we are assigning to the annotated step
   Class<?> value();

   // The name of the validate method in the view class
   String validateMethodName() default "validate";
}