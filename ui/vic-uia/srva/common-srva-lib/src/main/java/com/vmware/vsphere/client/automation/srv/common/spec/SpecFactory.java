/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import org.apache.commons.lang.RandomStringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.annotations.Test;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.SrvLocalizationUtil;
import com.vmware.client.automation.util.SsoUtil;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Spec factory class.
 */
public class SpecFactory {

   private static final int RANDOM_STRING_LENGTH = 5;
   private static final Logger _logger = LoggerFactory.getLogger(SpecFactory.class);

   private SpecFactory() {
   }

   /**
    * Generate spec with provided name and associated parent entity.
    *
    * @param specClass
    *           the class of the spec that will be generated
    * @param name
    *           the name of the spec that will be used
    * @param parent
    *           parent object that will be associated with the object spec
    * @return generated object spec
    */
   static public <T extends ManagedEntitySpec> T getSpec(Class<T> specClass,
         String name, ManagedEntitySpec parent) {

      T spec;
      try {
         spec = specClass.newInstance();
      } catch (Exception e) {
         String errorMessage =
               String.format(
                     "Unable to instantiate spec of type %s",
                     specClass.getName());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage);
      }

      spec.name.set(Strings.isNullOrEmpty(name) ? buildUniqueSpecName(specClass)
            : name);

      if (parent != null) {
         spec.parent.set(parent);
      }
      // Set service spec
      // TODO rkovachev: remove it once moved new test workflow
      ServiceSpec serviceSpec = SsoUtil.getServiceSpec(null);
      spec.service.set(serviceSpec);

      return spec;
   }

   /**
    * Generate spec with generated unique name and associated parent entity. The
    * name of the spec consists of object name and randomly generated string.
    * Note that the method doesn't check the testbed for existing objects so
    * there is minor chance that the returned name is not unique.
    *
    * @param specClass
    *           the class of the spec that will be generated
    * @param parent
    *           parent object that will be associated with the object spec
    * @return generated object spec
    */
   static public <T extends ManagedEntitySpec> T getSpec(Class<T> specClass,
         ManagedEntitySpec parent) {
      return getSpec(specClass, null, parent);
   }

   /**
    * Generate spec with generated unique name. The parent entity is not set.
    * The name of the spec consists of object name and randomly generated
    * string. Note that the method doesn't check the testbed for existing
    * objects so there is minor chance that the returned name is not unique.
    *
    * @deprecated - use factory method with parent parameter instead.
    *
    * @param specClass
    *           the class of the spec that will be generated
    * @return generated object spec
    */
   @Deprecated
   static public <T extends ManagedEntitySpec> T getSpec(Class<T> specClass) {
      return getSpec(specClass, null);
   }

   /**
    * Generates unique name to be used inside a object spec.
    *
    * @param specClass
    *           the class of the spec name that will be generated
    * @return generated specification name
    */
   public static String buildUniqueSpecName(Class<? extends EntitySpec> specClass) {
      // Extract test name and use it as prefix
      String testNamePrefix = extractTestName();
      return String.format(
            "%s%s%s",
            testNamePrefix,
            specClass.getSimpleName().replace("Spec", ""),
            RandomStringUtils.randomAlphanumeric(RANDOM_STRING_LENGTH));
   }

   /**
    * Generates unique description for specs.
    *
    * @return the generated unique description
    */
   public static String buildUniqueDesc() {
      return String.format(
            "Description_%s",
            RandomStringUtils.randomAlphanumeric(RANDOM_STRING_LENGTH));
   }

   /**
    * Gets user spec.
    * The generated user spec has random name and password.
    *
    * @return     generated user spec
    */
   public static UserSpec getUserSpec() {
      UserSpec user = new UserSpec();
      user.username.set(
            RandomStringUtils.randomAlphabetic(10) + "@" +
                  SrvLocalizationUtil.getLocalizedString("user.domain")
            );
      user.password.set("Admin!23" + RandomStringUtils.randomAlphanumeric(5));
      // Set service spec
      // TODO rkovachev: remove it once moved new test workflow
      ServiceSpec serviceSpec = SsoUtil.getServiceSpec(null);
      user.service.set(serviceSpec);
      return user;
   }

   /**
    * Browse the stack trace and extract the test name from it.
    * The method is used to prefix the entities created by the SpecFactory to map them to the respective tests.
    * If the test name is not found(like the case when the code is invoked by a provider) empty string is the result.
    * @return name of the test or empty string in case the code is not invoked in a test work flow.
    */
   private static String extractTestName() {
      String testName = "";
      // Test class contains execute method
      String baseWorkflowMethodName = "execute";

      StackTraceElement[] stackTraceElements = Thread.currentThread().getStackTrace();
      for (StackTraceElement stackTraceElement : stackTraceElements) {
         if(stackTraceElement.getMethodName().equals(baseWorkflowMethodName)) {
            try {
               Class<?> cls = Class.forName(stackTraceElement.getClassName());
               // It could be test or provider - check for testNG annotation
               if(cls.getMethod(baseWorkflowMethodName).getAnnotation(Test.class) != null) {
                  testName = cls.getSimpleName();
                  break;
               }
            } catch (ReflectiveOperationException e) {
               // Analyzing the stack trace throws exception - return empty string as a result
               _logger.warn("Reflective exception thrown analyzing the main thread stack trace: " + e.getMessage());
               break;
            }
         }
      }
      return testName;
   }
}