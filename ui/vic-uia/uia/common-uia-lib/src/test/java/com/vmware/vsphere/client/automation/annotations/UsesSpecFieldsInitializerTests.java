package com.vmware.vsphere.client.automation.annotations;

import java.util.Arrays;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import junit.framework.Assert;

import org.junit.Test;

import com.google.common.base.Supplier;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpecFieldsInitializer;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpecFieldsInitializer.UsesSpecInitializationException;
import com.vmware.vsphere.client.automation.common.annotations.MockedSpec;

public class UsesSpecFieldsInitializerTests {

   @Test
   public void initializePrivateField() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());

      Supplier<MockedSpec> rawSpecAnnotatedInstance = new Supplier<MockedSpec>() {
         @UsesSpec
         private MockedSpec privateSpecField;

         @Override
         public MockedSpec get() {
            return privateSpecField;
         }
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      initializer.initializeFields(workflowSpec);

      Assert.assertEquals("Private field is set",
            workflowSpec.get(MockedSpec.class), rawSpecAnnotatedInstance.get());
   }

   @Test
   public void initializeProtectedArrayField() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());
      workflowSpec.add(new MockedSpec());

      Supplier<MockedSpec[]> rawSpecAnnotatedInstance = new Supplier<MockedSpec[]>() {
         @UsesSpec
         protected MockedSpec[] protectedArraySpecField;

         @Override
         public MockedSpec[] get() {
            return protectedArraySpecField;
         }
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      initializer.initializeFields(workflowSpec);

      Assert.assertTrue("Protected field is set", Arrays.equals(
            workflowSpec.getAll(MockedSpec.class).toArray(new MockedSpec[0]),
            rawSpecAnnotatedInstance.get()));
   }

   @Test
   public void initializePublicListField() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());
      workflowSpec.add(new MockedSpec());

      Supplier<List<MockedSpec>> rawSpecAnnotatedInstance = new Supplier<List<MockedSpec>>() {
         @UsesSpec
         public List<MockedSpec> publicListSpecField;

         @Override
         public List<MockedSpec> get() {
            return publicListSpecField;
         }
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      initializer.initializeFields(workflowSpec);

      Set<MockedSpec> expectedFieldSpecSet = new HashSet<MockedSpec>();
      expectedFieldSpecSet.addAll(workflowSpec.getAll(MockedSpec.class));

      Set<MockedSpec> actualPrivateFieldSpecSet = new HashSet<MockedSpec>();
      actualPrivateFieldSpecSet.addAll(rawSpecAnnotatedInstance.get());

      Assert.assertTrue("Public field is set",
            expectedFieldSpecSet.equals(actualPrivateFieldSpecSet));

   }

   @Test
   public void initializePublicIterableField() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());
      workflowSpec.add(new MockedSpec());

      Supplier<Iterable<MockedSpec>> rawSpecAnnotatedInstance = new Supplier<Iterable<MockedSpec>>() {
         @UsesSpec
         public Iterable<MockedSpec> publicField;

         @Override
         public Iterable<MockedSpec> get() {
            return publicField;
         }
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      initializer.initializeFields(workflowSpec);

      Set<MockedSpec> expectedFieldSpecSet = new HashSet<MockedSpec>();
      expectedFieldSpecSet.addAll(workflowSpec.getAll(MockedSpec.class));

      Set<MockedSpec> actualPrivateFieldSpecSet = new HashSet<MockedSpec>();
      for (MockedSpec mockedSpec : rawSpecAnnotatedInstance.get()) {
         actualPrivateFieldSpecSet.add(mockedSpec);
      }

      Assert.assertTrue("Public field is set",
            expectedFieldSpecSet.equals(actualPrivateFieldSpecSet));

   }

   @Test(expected = UsesSpecInitializationException.class)
   public void annotationIsNotAllowedForNonSpecRawFields() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());
      workflowSpec.add(new MockedSpec());

      Object rawSpecAnnotatedInstance = new Object() {
         @UsesSpec
         public int thisSpecAnnotationIsWrong;
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      try {
         initializer.initializeFields(workflowSpec);
      } catch (UsesSpecInitializationException e) {
         Assert.assertEquals(e.getInitializationMessage(),
               "The field should inherit from PropertyBox or should be Array or List<> of classes inherited from property box.");

         throw e;
      }
   }

   @Test(expected = UsesSpecInitializationException.class)
   public void annotationIsNotAllowedForNonSpecArrayFields() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());
      workflowSpec.add(new MockedSpec());

      Object rawSpecAnnotatedInstance = new Object() {
         @UsesSpec
         public int[] thisSpecAnnotationIsWrong;
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      try {
         initializer.initializeFields(workflowSpec);
      } catch (UsesSpecInitializationException e) {
         Assert.assertEquals(e.getInitializationMessage(),
               "The type of the array elements should inherit PropertyBox");

         throw e;
      }
   }

   @Test(expected = UsesSpecInitializationException.class)
   public void annotationIsNotAllowedForNonSpecListFields() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());
      workflowSpec.add(new MockedSpec());

      Object rawSpecAnnotatedInstance = new Object() {
         @UsesSpec
         public List<Integer> thisSpecAnnotationIsWrong;
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      try {
         initializer.initializeFields(workflowSpec);
      } catch (UsesSpecInitializationException e) {
         Assert.assertEquals(e.getInitializationMessage(),
               "The type of the list elements should inherit from PropertyBox");

         throw e;
      }
   }

   @Test(expected = UsesSpecInitializationException.class)
   public void errorOnSigleRequestWhenSeveralSpecsArePublished() {
      WorkflowSpec workflowSpec = new WorkflowSpec();
      workflowSpec.add(new MockedSpec());
      workflowSpec.add(new MockedSpec());

      Object rawSpecAnnotatedInstance = new Object() {
         @UsesSpec
         private MockedSpec privateSpecField;
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      try {
         initializer.initializeFields(workflowSpec);
      } catch (UsesSpecInitializationException e) {
         Assert.assertEquals(
               "Several specs were found for field. Declare the field as array or List<T> (super of List<T>) to get all of them or filter the workflow to get single.",
               e.getInitializationMessage());

         throw e;
      }
   }

   @Test(expected = UsesSpecInitializationException.class)
   public void errorOnSigleRequestWhenNoSpecsArePublished() {
      WorkflowSpec workflowSpec = new WorkflowSpec();

      Object rawSpecAnnotatedInstance = new Object() {
         @UsesSpec
         private MockedSpec privateSpecField;
      };

      UsesSpecFieldsInitializer initializer = new UsesSpecFieldsInitializer(
            rawSpecAnnotatedInstance);
      try {
         initializer.initializeFields(workflowSpec);
      } catch (UsesSpecInitializationException e) {
         Assert.assertEquals(e.getInitializationMessage(),
               "There isn't spec with the corresponding type.");

         throw e;
      }
   }
}
