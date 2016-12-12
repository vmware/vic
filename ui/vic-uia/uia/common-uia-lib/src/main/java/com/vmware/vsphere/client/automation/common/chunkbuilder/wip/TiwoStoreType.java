/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.chunkbuilder.wip;

import com.vmware.client.automation.common.step.LoginStep;
import com.vmware.client.automation.common.step.LogoutStep;
import com.vmware.client.automation.common.step.LogoutWithWaitStep;
import com.vmware.client.automation.testbed.TestBed;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.hsua.common.datamodel.PropertyBoxLinks;
import com.vmware.vsphere.client.automation.common.chunkbuilder.TestChunkBuilder;
import com.vmware.vsphere.client.automation.common.step.GlobalRefreshStep;

/**
 * Specifies all available methods for storing TIWO.
 * Only two are available: logout/login and global refresh.
 *
 * If LOGOUT method is used a UserSpec has to be set in the test.
 */
public enum TiwoStoreType implements TestChunkBuilder {

    /**
     * This method of persisting TIWO/WIP objects to XML
     * requires logout operation. Steps that are performed are:
     * * logout
     * * login again with the user used for the test up to now
     */
    LOGOUT {
        @Override
        public void appendSpecs(TestBed testBed, PropertyBoxLinks links) {
            // no specs needed - expects that userSpec is provider by the test
        }

        @Override
        public void appendPrereqSteps(WorkflowComposition composition) {
            // no prerequisites required
        }

        @Override
        public void appendTestSteps(WorkflowComposition composition) {
            composition.appendStep(new LogoutStep());
            composition.appendStep(new LoginStep());
        }
    },

    /**
     * This method of persisting TIWO/WIP objects to XML
     * requires logout operation. Steps that are performed are:
     * * logout
     * * login again with the user used for the test up to now
     */
    LOGOUT_WITH_WAIT {
        @Override
        public void appendSpecs(TestBed testBed, PropertyBoxLinks links) {
            // no specs needed - expects that userSpec is provider by the test
        }

        @Override
        public void appendPrereqSteps(WorkflowComposition composition) {
            // no prerequisites required
        }

        @Override
        public void appendTestSteps(WorkflowComposition composition) {
            composition.appendStep(new LogoutWithWaitStep());
            composition.appendStep(new LoginStep());
        }
    },

    /**
     * This method of persisting TIWO/WIP objects to XML
     * requires only refreshing the NGC UI by clicking
     * the general refresh button.
     */
    GLOBAL_REFRESH {
        @Override
        public void appendSpecs(TestBed testBed, PropertyBoxLinks links) {
            // no specification are required
        }

        @Override
        public void appendPrereqSteps(WorkflowComposition composition) {
            // no prerequisites are required
        }

        @Override
        public void appendTestSteps(WorkflowComposition composition) {
            composition.appendStep(new GlobalRefreshStep());
        }
    },

    /**
     * This method of persisting TIWO/WIP objects to XML
     * requires nothing to be done between minimizing and
     * maximizing.
     */
    NONE {
        @Override
        public void appendSpecs(TestBed testBed, PropertyBoxLinks links) {
            // no specification are required
        }

        @Override
        public void appendPrereqSteps(WorkflowComposition composition) {
            // no prerequisites are required
        }

        @Override
        public void appendTestSteps(WorkflowComposition composition) {
            // not test steps required
        }
    }
}
