/**
 *
 */
package com.vmware.suitaf;

import static com.vmware.suitaf.SUITA.Factory.apl;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.SpecialStateHandler;
import com.vmware.suitaf.apl.SpecialStates;
import com.vmware.suitaf.apl.sele.SeleAPLImpl;

/**
 * This is the sub-tool {@link SubToolBrowser} with the following specs:
 * <li> <b>Function Type:</b> ACTIONS simulated at the host BROWSER
 * <li> <b>Description:</b> Execution of browser commands like open-url,
 * prev-page and next-page using the browser's JavaScript engine.
 * <li> <b>Based on APL:</b>
 * {@link AutomationPlatformLink#openUrl(String)}
 * {@link AutomationPlatformLink#goBack()}
 * {@link AutomationPlatformLink#goForward()}
 * <li> <b>Auxiliary SubTools:</b>
 * {@link SubToolSpecialState#getHandler(SpecialStates)},
 * {@link SubToolAudit#aplFailure(Throwable)}
 */
public class SubToolBrowser extends BaseSubTool {
    public SubToolBrowser(UIAutomationTool uiAutomationTool) {
        super(uiAutomationTool);
    }

    /**
     * Action Method that opens a web page in the browser requested through an
     * URL property.
     * @param urlProperty - the URL property of the page to be open
     */
    public void open(DataProperty<String> urlProperty) {
        open(urlProperty.get());
    }
    /**
     * Action Method that opens a web page in the browser requested through an
     * URL string.
     * @param urlString - the URL string of the page to be open
     */
    public void open(String urlString) {
        try {
            apl().openUrl(urlString);
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }

        // Handle a possible special state
        SpecialStateHandler ssh = ui.specialState.getHandler(
                SpecialStates.WIN_IE_CERT_ERROR_OVERRIDE_LINK);
        // Check if special state is detected
        int retryCount = 3;
        while (ssh.stateRecognize()) {
            if(retryCount > 0) {
               ssh.stateHandle();
               retryCount--;
            } else {
               ui.audit.aplFailure(new RuntimeException("Unable to hanlde untrusted browser certificate"));
            }
        }
    }

    /**
     * Action method that refreshes the current page.
     */
    public void refresh() throws InterruptedException {
       ((SeleAPLImpl) SUITA.Factory.apl()).getWebDriver().navigate().refresh();

       ui.condition.isFound("mainControlBar/loggedInUser").await(
             SUITA.Environment.getPageLoadTimeout() * 4);
    }

    /**
     * Action Method that commands the browser to open the previous page in
     * the page-open history
     */
    public void back() {
        try {
            apl().goBack();
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }
    }

    /**
     * Action Method that commands the browser to open the next page in the
     * page-open history
     */
    public void forward() {
        try {
            apl().goForward();
        } catch (Exception e) {
            ui.audit.aplFailure(e);
        }
    }
}