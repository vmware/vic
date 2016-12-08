/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.AbstractProperty;
import org.apache.commons.collections4.CollectionUtils;

import java.util.List;

/**
 * Utilities for making assertions about test specs.
 */
public class TestSpecValidator {

   /**
    * Ensures that the supplied argument is not null. Otherwise throws an
    * exception with the provided message.
    *
    * @param argument
    * @param message
    */
   public static void ensureNotNull(Object argument, String message) {
      ensureTrue(argument != null, message);
   }

   /**
    * Ensures that the supplied argument is not null or empty. Otherwise throws an
    * exception with the provided message.
    *
    * @param argument
    * @param message
    */
   public static void ensureNotBlank(AbstractProperty<String> argument, String message) {
      ensureTrue(!Strings.isNullOrEmpty(argument.get()), message);
   }

   /**
    * Ensures that the supplied List of argument is not empty. Otherwise throws an
    * exception with the provided message.
    *
    * @param argument
    * @param message
    */
   public static void ensureNotEmpty(List<? extends BaseSpec> argument, String
         message) {
      ensureTrue(!CollectionUtils.isEmpty(argument), message);
   }

   /**
    * Ensures that the supplied argument is assigned a value. Otherwise throws
    * an exception with the provided message.
    *
    * @param argument
    * @param message
    */
   public static void ensureAssigned(AbstractProperty<?> argument,
         String message) {
      ensureTrue(argument.isAssigned(), message);
   }

   /**
    * Ensures that the supplied argument's value is true. Otherwise throws an
    * exception with the provided message.
    *
    * @param argument
    * @param message
    */
   private static void ensureTrue(boolean argument, String message) {
      if (!argument) {
         throw new IllegalArgumentException(message);
      }
   }
}
