package com.vmware.vsphere.client.automation.storage.lib.core.tests;

import java.util.Arrays;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import org.junit.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.test.TestWorkflowStep;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;

import junit.framework.Assert;

public class WorkflowProcessorTest {

   private static class MockStep implements TestWorkflowStep {

      @Override
      public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      }

      @Override
      public void execute() throws Exception {
      }

      @Override
      public void clean() throws Exception {
      }

      @Override
      public List<RuntimeException> getFailedValidations() {
         return null;
      }

      @Override
      public void setStepTestScope(TestScope stepScope) {
      }

      @Override
      public TestScope getStepTestScope() {
         return null;
      }

      @Override
      public void logErrorInfo() {
      }
   }

   private static class BaseMockSpec extends BaseSpec {
   }

   private static class DerivedMockSpec extends BaseMockSpec {
   }

   private static class TestSpecData {
      public BaseSpec baseSpec = new BaseSpec();
      public BaseSpec baseMockSpec1 = new BaseMockSpec();
      public BaseSpec baseMockSpec2 = new BaseMockSpec();
      public BaseSpec baseMockSpec3 = new BaseMockSpec();
      public BaseSpec derivedMockSpec1 = new DerivedMockSpec();
      public BaseSpec derivedMockSpec2 = new DerivedMockSpec();
      public BaseSpec derivedMockSpec3 = new DerivedMockSpec();
   }

   @Test
   public void whitelistDerived() {
      WorkflowSpec containerSpec = new WorkflowSpec();
      TestSpecData testData = fillContainerSpec(containerSpec);

      WorkflowProcessor processor = new WorkflowProcessor(containerSpec);
      processor.append("SomeTitle", new MockStep(), testData.derivedMockSpec2,
            testData.derivedMockSpec3);

      // All spec should have same single tag except derivedMockSpec1
      List<String> specTags = testData.baseSpec.tag.getAll();
      Assert.assertEquals("Spec should have only one tag", 1, specTags.size());
      String tagName = specTags.get(0);

      verifyTags(testData.baseMockSpec1, tagName);
      verifyTags(testData.baseMockSpec2, tagName);
      verifyTags(testData.baseMockSpec3, tagName);

      verifyTags(testData.derivedMockSpec1); // No tags
      verifyTags(testData.derivedMockSpec2, tagName);
      verifyTags(testData.derivedMockSpec3, tagName);
   }

   @Test
   public void whitelistBase() {
      WorkflowSpec containerSpec = new WorkflowSpec();
      TestSpecData testData = fillContainerSpec(containerSpec);

      WorkflowProcessor processor = new WorkflowProcessor(containerSpec);
      processor.append("SomeTitle", new MockStep(), testData.baseMockSpec2,
            testData.baseMockSpec3);

      // Only the root and the specs in the white list should be tagged
      List<String> specTags = testData.baseSpec.tag.getAll();
      Assert.assertEquals("Spec should have only one tag", 1, specTags.size());
      String tagName = specTags.get(0);

      verifyTags(testData.baseMockSpec1); // No tags
      verifyTags(testData.baseMockSpec2, tagName);
      verifyTags(testData.baseMockSpec3, tagName);

      verifyTags(testData.derivedMockSpec1); // No tags
      verifyTags(testData.derivedMockSpec2); // No tags
      verifyTags(testData.derivedMockSpec3); // No tags
   }

   @Test
   public void whitelistRoot() {
      WorkflowSpec containerSpec = new WorkflowSpec();
      TestSpecData testData = fillContainerSpec(containerSpec);

      WorkflowProcessor processor = new WorkflowProcessor(containerSpec);
      processor.append("SomeTitle", new MockStep(), testData.baseSpec);

      // Only the root spec should be tagged
      List<String> specTags = testData.baseSpec.tag.getAll();
      Assert.assertEquals("Spec should have only one tag", 1, specTags.size());

      verifyTags(testData.baseMockSpec1); // No tags
      verifyTags(testData.baseMockSpec2); // No tags
      verifyTags(testData.baseMockSpec3); // No tags

      verifyTags(testData.derivedMockSpec1); // No tags
      verifyTags(testData.derivedMockSpec2); // No tags
      verifyTags(testData.derivedMockSpec3); // No tags
   }

   @Test
   public void existingTagsArePreserved() {
      String[] tagsArray = new String[] { "1", "2" };
      List<String> tagsList = Arrays.asList(tagsArray);

      WorkflowSpec containerSpec = new WorkflowSpec();
      BaseSpec derivedMockSpec1 = new DerivedMockSpec();
      derivedMockSpec1.tag.set(tagsList);
      BaseSpec derivedMockSpec2 = new DerivedMockSpec();
      derivedMockSpec2.tag.set(tagsList);
      containerSpec.add(derivedMockSpec1, derivedMockSpec2);

      WorkflowProcessor processor = new WorkflowProcessor(containerSpec);
      processor.append("SomeTitle", new MockStep(), derivedMockSpec1);

      // Only the root spec should be tagged
      List<String> specTags = derivedMockSpec1.tag.getAll();
      Assert.assertEquals("Spec preserved initial tags", true,
            specTags.containsAll(tagsList));

      verifyTags(derivedMockSpec2, tagsArray);
   }

   @Test
   public void composeSteps() {
      WorkflowSpec containerSpec = new WorkflowSpec();
      TestSpecData testData = fillContainerSpec(containerSpec);

      String step1Title = "Title1";
      String step2Title = "Title2";
      MockStep step1 = new MockStep();
      MockStep step2 = new MockStep();
      WorkflowProcessor processor = new WorkflowProcessor(containerSpec);
      processor.append(step1Title, step1, testData.baseMockSpec2,
            testData.baseMockSpec3);
      processor.append(step2Title, step2);

      WorkflowStepsSequence<TestWorkflowStepContext> flow = new WorkflowStepsSequence<>();
      processor.compose(flow);

      List<TestWorkflowStepContext> steps = flow.getAllSteps();
      TestWorkflowStepContext context = steps.get(0);
      Assert.assertEquals("Title is correct", step1Title, context.getTitle());
      Assert.assertEquals("Step is correct", step1, context.getStep());
      Assert.assertEquals("Step has one tag", 1,
            context.getSpecTagsFilter().length);

      context = steps.get(1);
      Assert.assertEquals("Title is correct", step2Title, context.getTitle());
      Assert.assertEquals("Step is correct", step2, context.getStep());
      Assert.assertEquals("Step doesn't have tags", null,
            context.getSpecTagsFilter());
   }

   private static void verifyTags(BaseSpec spec, String... tags) {
      Set<String> specTags = new HashSet<>(spec.tag.getAll());
      if (tags == null) {
         Assert.assertEquals("Spec doesn't have any tags", 0, specTags.size());
      } else {
         Assert.assertEquals("Spec has correct tags",
               new HashSet<String>(Arrays.asList(tags)), specTags);
      }
   }

   private static TestSpecData fillContainerSpec(BaseSpec containerSpec) {
      TestSpecData data = new TestSpecData();
      containerSpec.add(data.baseSpec, data.baseMockSpec1, data.baseMockSpec2,
            data.baseMockSpec3, data.derivedMockSpec1, data.derivedMockSpec2,
            data.derivedMockSpec3);
      return data;
   }
}
