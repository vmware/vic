/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.suitaf.apl.sele;

import java.net.MalformedURLException;
import java.net.URL;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import org.openqa.selenium.UnexpectedAlertBehaviour;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.firefox.FirefoxDriver;
import org.openqa.selenium.firefox.FirefoxProfile;
import org.openqa.selenium.ie.InternetExplorerDriver;
import org.openqa.selenium.remote.CapabilityType;
import org.openqa.selenium.remote.DesiredCapabilities;
import org.openqa.selenium.remote.RemoteWebDriver;

import com.google.common.base.Strings;
import com.vmware.suitaf.SUITA;

/**
 * This class initializes the connection to a particular WebDriver
 */
public class InitializeWebDriver {
    private WebDriver driver;

   /**
    * Initializes web driver
    *
    * @param wdtype
    *           - WebDriver type. Available options are Firefox (FF, RFF), Internet Explorer (IE,RIE), Google Chrome
    *           (GC, RGC), Safari (RSF)
    * @param hubAddress
    *           - Selenium RC / Grid IP Address
    * @param browserArgs
    *           - Command line arguments which are passed to the browser at browser launch time. Currently this is taken
    *           into account only for Google Chrome
    * @throws MalformedURLException
    */
    public InitializeWebDriver(WebDriverType wdtype, String hubAddress, String browserArgs) throws MalformedURLException {

        DesiredCapabilities capabilities = null;
        URL connectToURL = null;
        switch (wdtype) {
        case FF:
            driver = new FirefoxDriver();
            break;
        case RFF:
            capabilities = DesiredCapabilities.firefox();
            connectToURL = new URL(hubAddress);
            // Currently this one doesn't close alerts, added for consistency
            capabilities.setCapability(CapabilityType.UNEXPECTED_ALERT_BEHAVIOUR, UnexpectedAlertBehaviour.ACCEPT);
            // Firefox Profile to allow the cip plug-in
            FirefoxProfile profile = new FirefoxProfile();
            // Enable the npapi based plugin - it is now deprecated and that setting will be obsolete once removed.
            // TODO: rkovachev - remove it once the testing does not run against NPAPI cip.
            profile.setPreference("plugin.state.npvmwareclientsupportplugin", 2);
            // The new CIP needs to enable the vmware-csd protocol handlers. The next two lines register the
            // vmware-csd protocol and silently(now warning confirmation) accepts running it.
            profile.setPreference("network.protocol-handler.external.vmware-csd", true);
            profile.setPreference("network.protocol-handler.warn-external.vmware-csd", false);

            capabilities.setCapability(FirefoxDriver.PROFILE, profile);
            driver = new RemoteWebDriver(connectToURL, capabilities);
            break;
        case IE:
            // IE browser zoom level must be set to 100% so that the native mouse events
            // can be set to the correct coordinates
            // The Protected Mode settings must be set for each zone to have the same value.
            // The value can be on or off, as long as it is the same for every zone.
            // To set the Protected Mode settings, choose 'Internet Options' from the Tools menu,
            // and click on the Security tab.
            // For each zone, there will be a check box at the bottom of the tab labeled
            // "Enable Protected Mode" - make sure for all zones it is checked or alternatively
            // for all zones it is unchecked"
            // You need to either add the path to the IEDriver in the Path variable or use for example below code
            // System.setProperty("webdriver.ie.driver", "PATH_TO\\IEDriverServer_32bit.exe");
            // Note that there is performance issue with the 64bit IEWebDriver, so use the 32bit instead
            driver = new InternetExplorerDriver();
            break;
        case RIE:
            // You need to add the path to the IEDriver in the Path variable on the Node machine
            // or use for example -Dwebdriver.ie.driver=PATH_TO\\IEDriverServer_32bit.exe
            // as a command line parameter when starting the node
            // Note that there is performance issue with the 64bit IEWebDriver, so use the 32bit instead
            capabilities = DesiredCapabilities.internetExplorer();
            // Set the capabilities to accept if they can the unexpected pop ups, and then throw an
            // UnexpectedAlertException which could be futher handled
            // Check https://code.google.com/p/selenium/wiki/DesiredCapabilities for more info
            capabilities.setCapability(CapabilityType.UNEXPECTED_ALERT_BEHAVIOUR, UnexpectedAlertBehaviour.ACCEPT);
            // Hovering is enabled by default, but we don't want it, because it causes some problems like: flickering,
            // sending mouse over events every second and override the automation events
            capabilities.setCapability(InternetExplorerDriver.ENABLE_PERSISTENT_HOVERING, false);
            connectToURL = new URL(hubAddress);
            driver = new RemoteWebDriver(connectToURL, capabilities);
            driver.manage().timeouts().pageLoadTimeout(
                  SUITA.Environment.getPageLoadTimeout(), TimeUnit.MILLISECONDS);
            break;
        case GH:
            // You need to either add the path to the ChromeDriver in the Path variable or use for example below code
            // System.setProperty("webdriver.chrome.driver", "PATH_TO\\chromedriver.exe");
            driver = new ChromeDriver();
            break;
        case RGH:
            // You need to add the path to the ChromeDriver in the Path variable on the Node machine
            // or use for example -Dwebdriver.chrome.driver=PATH_TO\\chromedriver.exe
            // as a command line parameter when starting the node
            capabilities = DesiredCapabilities.chrome();
            // Currently this one doesn't close alerts, added for consistency
            capabilities.setCapability(CapabilityType.UNEXPECTED_ALERT_BEHAVIOUR, UnexpectedAlertBehaviour.ACCEPT);
            // Allow CIP plug-in
            ChromeOptions options = new ChromeOptions();
            // Enable the npapi based plugin - it is now deprecated and that setting will be obsolete once removed.
            // TODO: rkovachev - remove it once the testing does not run against NPAPI cip.
            options.addArguments("always-authorize-plugins=true");
            // test-type parameter removes the yellow banner notifying about the unsupported command line flag
            // "--ignore-certificate-errors"
            options.addArguments("test-type");

            if (!Strings.isNullOrEmpty(browserArgs)) {
                // Append provided command line arguments
                // List of supported Google Chrome arguments: http://peter.sh/experiments/chromium-command-line-switches/
                String[] args = browserArgs.split(",");
                for (String arg : args) {
                    options.addArguments(arg.trim());
                }
            }

            // The new CIP needs to enable the vmware-csd protocol handlers. The next two lines register the
            // vmware-csd protocol and silently(now warning confirmation) accepts running it.
            Map<String, Boolean> localState = new HashMap<String, Boolean>();
            localState.put("protocol_handler.excluded_schemes.vmware-csd", false);
            options.setExperimentalOption("localState", localState);

            capabilities.setCapability(ChromeOptions.CAPABILITY, options);
            connectToURL = new URL(hubAddress);
            driver = new RemoteWebDriver(connectToURL, capabilities);
            break;
        case RSF:
            capabilities = DesiredCapabilities.safari();
            capabilities.setCapability(CapabilityType.UNEXPECTED_ALERT_BEHAVIOUR, UnexpectedAlertBehaviour.ACCEPT);
            connectToURL = new URL(hubAddress);
            driver = new RemoteWebDriver(connectToURL, capabilities);
            break;
        case SF:
            throw new IllegalArgumentException("Safari is not yet supported!");
        default:
            throw new IllegalArgumentException("The provided browser is not yet supported!");
        }

        // Set page load timeout
        driver.manage().timeouts().pageLoadTimeout(
              SUITA.Environment.getPageLoadTimeout(), TimeUnit.MILLISECONDS);
    }

    /**
     * Get the initialized web driver.
     *
     * @return
     */
    public WebDriver getDriver() {
        return driver;
    }
}