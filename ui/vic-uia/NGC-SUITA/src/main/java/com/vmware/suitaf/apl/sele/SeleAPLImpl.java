package com.vmware.suitaf.apl.sele;

import static com.vmware.suitaf.util.CommonUtils.smartEqual;

import java.awt.Point;
import java.io.File;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import org.apache.commons.io.FileUtils;
import org.openqa.selenium.By;
import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.Keys;
import org.openqa.selenium.OutputType;
import org.openqa.selenium.TakesScreenshot;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.interactions.Actions;
import org.openqa.selenium.remote.Augmenter;
import org.openqa.selenium.server.commands.NativeCommand;
import org.openqa.selenium.support.ui.Select;

import com.thoughtworks.selenium.FlashSelenium;
import com.thoughtworks.selenium.Selenium;
import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.CheckBox;
import com.vmware.flexui.componentframework.controls.mx.ComboBase;
import com.vmware.flexui.componentframework.controls.mx.ComboBox;
import com.vmware.flexui.componentframework.controls.mx.DataGrid;
import com.vmware.flexui.componentframework.controls.mx.ListBase;
import com.vmware.flexui.componentframework.controls.mx.NavBar;
import com.vmware.flexui.componentframework.controls.mx.NumericStepper;
import com.vmware.flexui.componentframework.controls.mx.RadioButton;
import com.vmware.flexui.componentframework.controls.mx.TextInput;
import com.vmware.flexui.componentframework.controls.mx.Tree;
import com.vmware.flexui.componentframework.controls.mx.custom.PermanentTabBar;
import com.vmware.flexui.componentframework.controls.mx.custom.RichComboBox;
import com.vmware.flexui.componentframework.controls.spark.SparkCheckBox;
import com.vmware.flexui.componentframework.controls.spark.SparkDropDownList;
import com.vmware.flexui.componentframework.controls.spark.SparkRadioButton;
import com.vmware.flexui.selenium.VMWareSelenium;
import com.vmware.suitaf.apl.AutomationPlatformLinkExt;
import com.vmware.suitaf.apl.Category;
import com.vmware.suitaf.apl.ComponentMatcher;
import com.vmware.suitaf.apl.HostLogger;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Key;
import com.vmware.suitaf.apl.MouseButtonAction;
import com.vmware.suitaf.apl.MouseComponentAction;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.apl.SpecialStateHandler;
import com.vmware.suitaf.apl.SpecialStates;
import com.vmware.suitaf.apl.WindowModes;
import com.vmware.suitaf.apl.sele.SeleHelper.KeyItem;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.suitaf.util.Condition;
import com.vmware.suitaf.util.FailureCauseHolder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SeleAPLImpl implements AutomationPlatformLinkExt {

   boolean useWebDriver = false;
   
   private Boolean isFlexApp = null;

   private static final String CSS_PREFIX = "css=";
   private static final String JS_PREFIX = "js=";

   private static final Logger _logger = LoggerFactory.getLogger(SeleAPLImpl.class);

   // Main operational resources
   protected SeleLinks seleLinks;
   protected SeleHelper seleHelper;
   protected SikuliXHelper sikuliXHelper;
   protected HostLogger hostLogger;

   @Override
   public void resetLink(
         Map<String, String> linkInitParams, HostLogger hostLogger) {
      // Force closure of the browser if such was left open
      if (this.seleLinks != null) {
         this.seleLinks.closeBrowserInstance();
         this.seleLinks = null;
         this.seleHelper = null;
         this.hostLogger = null;
      }

      // Refreshing configuration parameters
      final String serverIP = linkInitParams.get(InitParams.PN_SERVER_IP);
      final String browser = linkInitParams.get(InitParams.PN_BROWSER);
      final String browserArgs = linkInitParams.get(InitParams.PN_BROWSER_ARGS);
      final String baseURL = linkInitParams.get(InitParams.PN_BASE_URL);

      if (serverIP == null || browser == null || baseURL == null) {
         // If any parameter is null - shut down
         return;
      }

      final String unparsedPort = linkInitParams.get(InitParams.PN_SERVER_PORT);
      int serverPort = Integer.parseInt(unparsedPort);
      this.hostLogger = hostLogger;
      this.useWebDriver = Boolean.parseBoolean(
            linkInitParams.get(InitParams.PN_USE_WEB_DRIVER));

      this.seleLinks = new SeleLinks(this, serverIP, serverPort, browser, browserArgs, baseURL);

      this.seleHelper = new SeleHelper(this);

      seleLinks.setPageLoadTimeout(linkInitParams.get(InitParams.PN_PAGE_LOAD_TIMEOUT));

      seleLinks.openBrowserInstance();

      final String targetHostName = SeleLinks.getWebDriverNode().getHostName();
      this.sikuliXHelper = new SikuliXHelper(targetHostName);
   }

   @Override
   public void attachLink(
         Map<String, String> linkInitParams, HostLogger hostLogger) {
      hostLogger.error(
            "Unsupported method <attachLink>! Redirecting to <resetLink>.");
      resetLink(linkInitParams, hostLogger);
   }

   @Override
   public void captureBitmap(String fileName) {
      if (seleLinks.selenium != null) {
         seleLinks.selenium.captureScreenshot(fileName);
      } else {
         // The screenshot won't work (produce black rectangle) for IE if there has not been made
         // GUI connection to the remote machine from the machine running the tests
         WebDriver augmentedDriver = new Augmenter().augment(seleLinks.driver);
         try {
            File screenshot =
                  ((TakesScreenshot) augmentedDriver).getScreenshotAs(OutputType.FILE);
            FileUtils.copyFile(screenshot, new File(fileName));
         } catch (Exception e) {
            e.printStackTrace();
         }
      }
   }

   protected final DirectID DESKTOP = DirectID.directid(
         this, SeleLinks.PATH_TO_DESKTOP);
   protected final DirectID BROWSER_APP = DirectID.directid(
         this, SeleLinks.PATH_TO_BROWSER_APP);
   protected final DirectID BROWSER_WIN = DirectID.directid(
         this, SeleLinks.PATH_TO_BROWSER_WIN);

   @Override
   public ComponentMatcher getComponentMatcher(IDGroup... idGroups) {
      return new SeleComponentID(this, idGroups);
   }

   @Override
   public int getExistingCount(ComponentMatcher componentID) {
      DirectID directId = getDirectId(componentID);
      return getExistingCount(directId);
   }

   private int getExistingCount(DirectID directId) {
      int count;
      if (directId.hasGraphicalID()) {
         final String graphicalID = directId.getGraphicalID();
         count = sikuliXHelper.getExistingCount(graphicalID);
      } else if (directId.isRawType()) {
         count = seleHelper.getExistingCount(directId);
      } else {
         if(isFlexApp()) {
            count = seleLinks.getExistingCount(directId);
         } else {
            count = findWebElements(directId).size();
         }
      }
      return count;
   }

   @Override
   public List<List<String>> getGridProperty(
         ComponentMatcher baseComponent, Property property) {
      // Retrieve a proxy-instance for the target component
      final DirectID directID = getDirectId(baseComponent);
      DisplayObject displayObject = seleLinks.find(directID);
      // Retrieve the list of property values for the requested property
      return seleHelper.getGridProperty(displayObject, property);
   }

   @Override
   public List<String> getListProperty(
         ComponentMatcher baseComponent, Property property) {
      // Retrieve a proxy-instance for the target component
      final DirectID directID = getDirectId(baseComponent);
      DisplayObject displayObject = seleLinks.find(directID);
      // Retrieve the list of property values for the requested property
      return seleHelper.getListProperty(displayObject, property);
   }

   @Override
   public String getSingleProperty(
         ComponentMatcher baseComponent, Property property) {
      final DirectID directID = getDirectId(baseComponent);
      return getSingleProperty(directID, property);
   }

   private boolean isFlexApp() {
      if(isFlexApp == null) {
         isFlexApp = getFlashSelenium() != null;
      }
      return isFlexApp;
   }

   /**
    * This method returns the value of a property for a specified object.
    * Depending on the data in the object, Sikuli or Selenium
    * implementation is used.
    *
    * @param baseComponent - the object which property value is to be returned.
    * @param property - the needed property of the object.
    *
    * @return - value of the needed property as a String.
    */
   private String getSingleProperty(
         DirectID baseComponent, Property property) {
      String propertyValue;

      if (baseComponent.hasGraphicalID()) {
         final String graphicalID = baseComponent.getGraphicalID();
         propertyValue = sikuliXHelper.getSingleProperty(graphicalID, property);
      } else if (baseComponent.isRawType()) {
         propertyValue = seleHelper.getRawProperty(baseComponent, property);
      } else {
         if(isFlexApp()) {
            DisplayObject displayObject = getDisplayObject(baseComponent);
            propertyValue = seleHelper.getSingleProperty(displayObject, property);
         } else {
            propertyValue = getElementProperty(baseComponent, property);
         }
      }
      return propertyValue;
   }

   private String getElementProperty(DirectID baseComponent, Property property) {
      WebElement webElement = findWebElement(baseComponent);
      // Mapping between the general properties and native property names
      switch (property) {
         case ID:
            return webElement.getAttribute("id");
         case ENABLED:
            return webElement.isEnabled() + "";
         case VALUE:
            return webElement.getText();
         case VISIBLE:
            return webElement.isDisplayed() + "";
         default:
            return null;
      }
   }

   private DisplayObject getDisplayObject(DirectID baseComponent) {
      // Retrieve a proxy-instance for the target component
      FailureCauseHolder fch = new FailureCauseHolder();
      DisplayObject displayObject = null;
      try {
         displayObject = seleLinks.find(baseComponent);
      } catch (Exception e) {
         fch.setCause(e);
         fch.escalateCause();
      }
      return displayObject;
   }

   private Integer getPropertyInt(DirectID baseComponent, Property property) {
      String value = getSingleProperty(baseComponent, property);
      if (value != null) {
         return Integer.valueOf(value);
      }
      return null;
   }

   // =====================================================================
   @Override
   public void executeJavaScript(String code) {
      seleLinks.selenium.getEval(code);
   }

   @Override
   public WindowModes getBrowserWindowMode() {
      throw new RuntimeException("Function not implemented.");
   }

   @Override
   public void setBrowserWindowMode(WindowModes mode) {
      throw new RuntimeException("Function not implemented.");
   }

   @Override
   public String getUrl() {
      return seleLinks.selenium.getLocation();
   }

   @Override
   public void goBack() {
      seleLinks.selenium.goBack();
   }

   @Override
   public void goForward() {
      ((VMWareSelenium) seleLinks.selenium).goForward();
   }

   @Override
   public void openUrl(String url) {
      seleLinks.openUrl(url);
   }

   @Override
   public void reloadPage() {
      seleLinks.selenium.refresh();
   }

   @Override
   public void mouseOnComponent(
         ComponentMatcher clickableComponent,
         MouseComponentAction action,
         Condition asserter) {
      DirectID directId = getDirectId(clickableComponent);

      if(!isFlexApp()) {
         WebElement element = findWebElement(directId);
         Actions actions = new Actions(getWebDriver());
         actions.moveToElement(element).click().perform();
         return;
      }

      // are we dealing with a raw component?
      if (directId.isRawType()) {
         this.seleHelper.clickRawComponent(directId);
         return;
      }

      // Retrieve a proxy-instance for the target directId
      UIComponent to = (UIComponent) seleLinks.find(directId);

      FailureCauseHolder ward = new FailureCauseHolder();
      int retry = 3;
      do {
         retry--;
         ward.setCause(null);

         try {
            switch (action) {
               case CLICK:
                  switch (retry) {
                     case 2:
                        to.leftMouseClick();
                        break;
                     case 1:
                        to.visibleMouseDownUp();
                        break;
                     case 0:
                        to.visibleClick();
                        break;
                  }
                  break;
               case RIGHT_CLICK:
                  to.rightClick();
                  break;
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }

         try {
            // Use the asserter to check if this action succeeded
            if (ward.isNull() && asserter != null) {
               asserter.verify();
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }
      } while (retry > 0 && !ward.isNull());

      ward.escalateCause();
   }

   @Override
   public void setFocus(ComponentMatcher focusableComponent) {
      final DirectID directId = getDirectId(focusableComponent);
      setFocus(directId);
   }

   protected void setFocus(DirectID focusableID) {
      if (isFlexApp()) {
         // Retrieve a proxy-instance for the target component
         UIComponent to = (UIComponent) seleLinks.find(focusableID);
         to.setFocus(1000);
      } else {
         WebElement element = findWebElement(focusableID);
         // Input tag need special handling for setting focus -
         // check Selenium documentation.
         if ("input".equals(element.getTagName())) {
            element.sendKeys("");
         } else {
            new Actions(getWebDriver()).moveToElement(element).perform();
         }
      }
   }

   @Override
   public void typeKeys(
         ComponentMatcher focusableComponent,
         int typeDelay,
         Object... keySequence) {

      if (seleLinks.driver() != null) {
         // NOTE: it is a quick and dirty implementation for the web
         // driver to be used to run all the BAT tests. support only the basic
         // key operation enter and tab.
         // TODO: rkovachev provide proper implementation.
         Object keyItem = keySequence[0];
         if (keyItem instanceof Key) {
            Keys keysToSend = null;
            Key key = (Key) keyItem;
            switch (key) {
               case ENTER:
                  keysToSend = Keys.ENTER;
                  break;
               case TAB:
                  keysToSend = Keys.TAB;
                  break;
               default:
                  throw new RuntimeException(
                        "The provided key operation is not implemented!");
            }

            seleLinks.driver().switchTo().activeElement().sendKeys(keysToSend);
         } else {
            throw new RuntimeException("Implement it!");
         }
         return;
      }

      if (focusableComponent == null) {
         typeKeys((DirectID) null, typeDelay, keySequence);
      } else {
         typeKeys(
               ((SeleComponentID) focusableComponent).mainID,
               typeDelay,
               keySequence);
      }
   }

   protected void typeKeys(
         DirectID typingCompID, int typeDelay, Object... keySequence) {
      // Parameters verification
      if (keySequence == null || keySequence.length == 0) {
         return;
      }

      if (typingCompID != null) {
         setFocus(typingCompID);
      }

      ArrayList<String> holdedKeys = new ArrayList<String>();

      for (KeyItem code : SeleHelper.getKeyItems(keySequence)) {
         if (code.isHoldKey) {
            seleLinks.selenium.keyDownNative(code.value);
            holdedKeys.add(0, code.value);
         } else if (code.isPasteText) {
            seleLinks.selenium.getEval(
                  NativeCommand.CLIPBOARD_PUT_TEXT.encode(code.value));
         } else {
            seleLinks.selenium.keyPressNative(code.value);
            if (holdedKeys.size() > 0) {
               for (String holded : holdedKeys) {
                  seleLinks.selenium.keyUpNative(holded);
               }
               holdedKeys.clear();
            }
            // Wait the delay between key types
            CommonUtils.sleep(typeDelay);
         }
      }
      if (holdedKeys.size() > 0) {
         for (String holded : holdedKeys) {
            seleLinks.selenium.keyUpNative(holded);
         }
         holdedKeys.clear();
      }
   }

   @Override
   public void setValue(
         ComponentMatcher componentId,
         String newValue,
         Condition asserter) {

      if(isFlexApp()) {
         setFlexValue(componentId, newValue, asserter);
      } else {
         setHtmlValue(componentId, newValue);
      }
   }

   /**
    * Find the first WebElement for the provided identifier.
    * @param id web control identifier
    * @return the first found WebElement or throws NoSuchElementException.
    */
   private WebElement findWebElement(DirectID id) {
      String JS_ELEMENT_RETRIEVAL_SCRIPT =
            "var resultingElement = %s; return resultingElement;";
      WebElement result = null;
      String htmlId = id.getHtmlID();
      if (htmlId.startsWith("//") || htmlId.startsWith("./") || htmlId.startsWith("../") || htmlId.equals(".")) {
         result = getWebDriver().findElement(By.xpath(htmlId));
      } else if (htmlId.startsWith(CSS_PREFIX)) {
         result = getWebDriver().findElement(By.cssSelector(htmlId.substring(CSS_PREFIX.length())));
      } else if (htmlId.startsWith(JS_PREFIX)) {
         JavascriptExecutor js = (JavascriptExecutor) getWebDriver();
         Object jsResult =
               js.executeScript(String.format(JS_ELEMENT_RETRIEVAL_SCRIPT, htmlId.substring(JS_PREFIX.length())));

         if (jsResult instanceof WebElement) {
            result = (WebElement) jsResult;
         }
      } else {
         result = getWebDriver().findElement(By.id(htmlId));
      }
      if(result != null && !result.isDisplayed()) {
         _logger.error("The list contains not visible controls: " + result.toString());
      }
      return result;
   }

   /**
    * Find all elements within the current page using the given identifier.
    * @param id control identifier.
    * @return a list of all WebElements, or an empty list if nothing matches
    */
   @SuppressWarnings("unchecked")
   private List<WebElement> findWebElements(DirectID id) {
      String JS_ELEMENT_RETRIEVAL_SCRIPT =
            "var resultingElement = %s; return resultingElement;";

      List<WebElement> result = new ArrayList<WebElement>();

      String htmlId = id.getHtmlID();
      if (htmlId.startsWith("//") || htmlId.startsWith("./") || htmlId.startsWith("../") || htmlId.equals(".")) {
         result = getWebDriver().findElements(By.xpath(htmlId));
      } else if (htmlId.startsWith(CSS_PREFIX)) {
         result = getWebDriver().findElements(By.cssSelector(htmlId.substring(CSS_PREFIX.length())));
      } else if (htmlId.startsWith(JS_PREFIX)) {
         JavascriptExecutor js = (JavascriptExecutor) getWebDriver();
         Object jsResult =
               js.executeScript(String.format(JS_ELEMENT_RETRIEVAL_SCRIPT, htmlId.substring(JS_PREFIX.length())));

         if (jsResult instanceof WebElement) {
            result.add((WebElement) jsResult);
         } else if (jsResult instanceof List<?>) {
            result = (List<WebElement>) jsResult;
         }
      } else {
         result = getWebDriver().findElements(By.id(htmlId));
      }
      for (WebElement webElement : result) {
         if(!webElement.isDisplayed()) {
            _logger.error("The list contains not visible controls: " + webElement.toString());
         }
      }
      return result;
   }

   private void setHtmlValue(ComponentMatcher componentId, String newValue) {
      // Type-cast the component ID to local ID implementation
      SeleComponentID localID = SeleComponentID.fromMatcher(componentId);
      WebElement webElement = findWebElement(localID.mainID);
      if(webElement.getTagName().equalsIgnoreCase("select")) {
         Select dropDown = new Select(webElement);
         dropDown.selectByVisibleText(newValue);
      } else {
         Actions action = new Actions(getWebDriver());
         action.sendKeys(webElement, newValue).build().perform();
      }
   }

   private void setFlexValue(ComponentMatcher componentId,
         String newValue,
         Condition asserter) {

      // Type-cast the component ID to local ID implementation
      SeleComponentID localID = SeleComponentID.fromMatcher(componentId);

      // Retrieve a proxy-instance for the target component
      DisplayObject to = seleLinks.find(localID.mainID);
      String currentValue = seleHelper.getSingleProperty(to, Property.VALUE);

      FailureCauseHolder ward = new FailureCauseHolder();
      int retry = 4;
      do {
         retry--;

         // Check if value already set
         if (smartEqual(currentValue, newValue)) {
            ward.setCause(null);
            break;
         }

         // Invoke setting the value to the component
         try {
            if (to instanceof SparkCheckBox) {
               try {
                  ((SparkCheckBox) to).checkUncheckCheckBox(newValue);
               } catch (AssertionError e) {
                  to = new CheckBox(
                        to.getUniqueId(), this.seleLinks.flashSelenium());
               }
            } else if (to instanceof CheckBox) {
               switch (retry) {
                  case 0:
                  case 2:
                     ((CheckBox) to).checkUncheckcheckbox(newValue);
                     break;
                  case 1:
                  case 3:
                     ((CheckBox) to).checkUncheckListCheckBox(newValue);
                     break;
               }
            } else if (to instanceof SparkRadioButton) {
               ((SparkRadioButton) to).visibleMouseDownUp();
            } else if (to instanceof RadioButton) {
               ((RadioButton) to).visibleMouseDownUp();
            } else if (to instanceof PermanentTabBar) {
               int tabsCount = ((PermanentTabBar) to).getNumTabs();
               for (int i = 0; i < tabsCount; i++) {
                  ((PermanentTabBar) to).selectTabAt(i + "", 0);
                  String selectedTabId =
                        ((PermanentTabBar) to).getChildIdAtIndex(i + "");
                  String tabName = seleHelper.safeGetProperty(
                        selectedTabId, true, "label");
                  if (tabName.equals(newValue)) {
                     break;
                  }
               }
            } else if (to instanceof NavBar) {
               ((NavBar) to).selectContextMenuItemByName(newValue);
            } else if (to instanceof TextInput) {
               ((TextInput) to).type(newValue);
            } else if (to instanceof NumericStepper) {
               ((NumericStepper) to).setValue(newValue);
            } else if (to instanceof Tree) {
               List<String> treeItems = seleHelper.getListProperty(
                     to, Property.VALUE_LIST);
               int valueIndex = treeItems.indexOf(newValue);
               if (valueIndex > -1) {
                  to.setProperty("selectedIndex", valueIndex + "");
                  List<String> pathSplit = Property.Convert.stringToTreePath(
                        newValue);

                  String treeLeafID = to.getUniqueId() + "/" +
                                      DirectID.SELE_IDATTRIB_AUTOMATIONNAME +
                                      "=" + pathSplit.get(pathSplit.size() - 1);
                  new UIComponent(
                        treeLeafID, seleLinks.flashSelenium()).click();
               }
            } else if (to instanceof ListBase) {
               switch (retry) {
                  case 1:
                     ((ListBase) to).selectItem(newValue);
                     break;
                  case 2:
                     ((ListBase) to).selectItemByField("text", newValue);
                     break;
                  case 3:
                     ((ListBase) to).selectItemByField("label", newValue);
                     break;
               }
            } else if (to instanceof RichComboBox) {
               ((RichComboBox) to).setValue(newValue);
            } else if (to instanceof ComboBox) {
               switch (retry) {
                  case 1:
                     ((ComboBox) to).select(newValue);
                     break;
                  case 2:
                     ((ComboBox) to).selectByLabel(newValue);
                     break;
                  case 3:
                     List<String> cbValues = getListProperty(
                           componentId, Property.VALUE_LIST);
                     if (cbValues != null && cbValues.contains(newValue)) {
                        ((ComboBox) to).selectItemByIndex(
                              cbValues.indexOf(newValue) + "");
                     }
                     break;
               }
            } else if (to instanceof SparkDropDownList) {
               SparkDropDownList ddList = ((SparkDropDownList) to);
               String indexToSelect = ddList.getIndexForValue(newValue);
               Integer index = new Integer(indexToSelect);
               if (index >= 0) {
                  ((SparkDropDownList) to).selectItemByIndex(indexToSelect);
               } else {
                  throw new RuntimeException(
                        "The item " + newValue +
                        " is not found in the drop down list!");
               }
            } else {
               String message = "Function <setValue> is not implemented" +
                                " for: " + to.getUniqueId();
               hostLogger.error(message);
               seleHelper.logFlexObject(to, 1, null, Category.ANY);
               ward.setCause(new RuntimeException(message));
               // Leave the loop early
               break;
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }

         try {

            currentValue = seleHelper.getSingleProperty(to, Property.VALUE);
            // Verify action through check of the current value
            if (ward.isNull() && !smartEqual(currentValue, newValue)) {
               ward.setCause(
                     new RuntimeException(
                           String.format(
                                 "Value setting to '%s' failed. Remains '%s'.",
                                 newValue,
                                 currentValue)));
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }

         try {
            if (ward.isNull() && asserter != null) {
               asserter.verify();
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }
      } while (retry > 0);

      ward.escalateCause();
   }

   @Override
   public void setValueByIndex(
         ComponentMatcher indexableComponent,
         int newIndex,
         Condition asserter) {
      // Type-cast the component ID to local ID implementation
      SeleComponentID localID = SeleComponentID.fromMatcher(indexableComponent);

      // Retrieve a proxy-instance for the target component
      DisplayObject to = seleLinks.find(localID.mainID);

      // Check if value already set
      int currentIndex = getPropertyInt(localID.mainID, Property.VALUE_INDEX);

      FailureCauseHolder ward = new FailureCauseHolder();
      int retry = 4;
      do {
         retry--;
         if (currentIndex == newIndex) {
            ward.setCause(null);
            break;
         }

         // Invoke setting the focus to the component
         try {
            // Complex workaround cases that are handled separately
            if (to instanceof DataGrid) {
               switch (retry) {
                  case 0:
                  case 2:
                     ((DataGrid) to).selectRow(newIndex + "");
                     break;
                  case 1:
                  case 3:
                     ((ListBase) to).selectItem(newIndex + "");
                     break;
               }
            } else if (to instanceof Tree) {
               List<String> treeItems = seleHelper.getListProperty(
                     to, Property.VALUE_LIST);

               to.setProperty("selectedIndex", newIndex + "");
               List<String> pathSplit = Property.Convert.stringToTreePath(
                     treeItems.get(newIndex));

               String treeLeafID = to.getUniqueId() + "/" +
                                   DirectID.SELE_IDATTRIB_AUTOMATIONNAME +
                                   "=" + pathSplit.get(pathSplit.size() - 1);
               new UIComponent(
                     treeLeafID, seleLinks.flashSelenium()).click();
            } else if (to instanceof ListBase) {
               ((ListBase) to).selectItem(newIndex + "");
            } else if (to instanceof ComboBase) {
               ((ComboBase) to).selectItemByIndex(newIndex + "");
            } else {
               String message = "Function <setValueByIndex> is not" +
                                " implemented for: " + to.getUniqueId();
               hostLogger.error(message);
               seleHelper.logFlexObject(to, 1, null, Category.ANY);
               throw new RuntimeException(message);
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }

         try {
            // Re-get the current index
            currentIndex = getPropertyInt(
                  localID.mainID, Property.VALUE_INDEX);

            // Verify action through check of the current value
            if (currentIndex != newIndex) {
               throw new RuntimeException(
                     String.format(
                           "Index setting to '%s' failed. Remains '%s'.",
                           newIndex,
                           currentIndex));
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }

         try {
            if (ward.isNull() && asserter != null) {
               asserter.verify();
            }
         } catch (Throwable e) {
            ward.setCause(e);
         }
      } while (retry > 0 && !ward.isNull());

      ward.escalateCause();
   }

   @Override
   public SpecialStateHandler getSpecialStateHandler(SpecialStates state) {
      return SeleSpecialStateHandlers.getHandler(state);
   }

   @Override
   public void mouseOnApplication(
         Condition asserter, Object... mouseActionsSequence) {

      FailureCauseHolder ward = new FailureCauseHolder();
      try {
         for (Object mouseAction : mouseActionsSequence) {
            if (mouseAction instanceof Long) {
               // Wait the delay between mouse actions
               CommonUtils.sleep((Long) mouseAction);
            } else if (mouseAction instanceof Point) {
               Point p = (Point) mouseAction;
               seleLinks.selenium.getEval(
                     NativeCommand.FLEXAPP_MOUSE_MOVE.encode(p.x, p.y));
            } else if (mouseAction instanceof MouseButtonAction) {
               MouseButtonAction action = (MouseButtonAction) mouseAction;
               switch (action) {
                  case LEFT_DOWN:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_DOWN.encode(1));
                     break;
                  case LEFT_UP:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_UP.encode(1));
                     break;
                  case LEFT_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_CLICK.encode(1));
                     break;
                  case LEFT_DOUBLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_DOUBLE_CLICK.encode(1));
                     break;
                  case MIDDLE_DOWN:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_DOWN.encode(2));
                     break;
                  case MIDDLE_UP:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_UP.encode(2));
                     break;
                  case MIDDLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_CLICK.encode(2));
                     break;
                  case MIDDLE_DOUBLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_DOUBLE_CLICK.encode(2));
                     break;
                  case RIGHT_DOWN:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_DOWN.encode(3));
                     break;
                  case RIGHT_UP:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_UP.encode(3));
                     break;
                  case RIGHT_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_CLICK.encode(3));
                     break;
                  case RIGHT_DOUBLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.FLEXAPP_MOUSE_DOUBLE_CLICK.encode(3));
                     break;
                  default:
                     throw new RuntimeException(
                           "Unsupported mouse button action: " + action);
               }
            }
         }
      } catch (Throwable e) {
         ward.setCause(e);
      }

      try {
         if (ward.isNull() && asserter != null) {
            asserter.verify();
         }
      } catch (Throwable e) {
         ward.setCause(e);
      }

      ward.escalateCause();
   }

   @Override
   public void mouseOnScreen(
         Condition asserter, Object... mouseActionsSequence) {

      FailureCauseHolder ward = new FailureCauseHolder();
      try {
         for (Object mouseAction : mouseActionsSequence) {
            if (mouseAction instanceof Long) {
               // Wait the delay between mouse actions
               CommonUtils.sleep((Long) mouseAction);
            } else if (mouseAction instanceof Point) {
               Point p = (Point) mouseAction;
               seleLinks.selenium.getEval(
                     NativeCommand.SCREEN_MOUSE_MOVE.encode(p.x, p.y));
            } else if (mouseAction instanceof MouseButtonAction) {
               MouseButtonAction action = (MouseButtonAction) mouseAction;
               switch (action) {
                  case LEFT_DOWN:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_DOWN.encode(1));
                     break;
                  case LEFT_UP:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_UP.encode(1));
                     break;
                  case LEFT_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_CLICK.encode(1));
                     break;
                  case LEFT_DOUBLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_DOUBLE_CLICK.encode(1));
                     break;
                  case MIDDLE_DOWN:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_DOWN.encode(2));
                     break;
                  case MIDDLE_UP:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_UP.encode(2));
                     break;
                  case MIDDLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_CLICK.encode(2));
                     break;
                  case MIDDLE_DOUBLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_DOUBLE_CLICK.encode(2));
                     break;
                  case RIGHT_DOWN:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_DOWN.encode(3));
                     break;
                  case RIGHT_UP:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_UP.encode(3));
                     break;
                  case RIGHT_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_CLICK.encode(3));
                     break;
                  case RIGHT_DOUBLE_CLICK:
                     seleLinks.selenium.getEval(
                           NativeCommand.SCREEN_MOUSE_DOUBLE_CLICK.encode(3));
                     break;
                  default:
                     throw new RuntimeException(
                           "Unsupported mouse button action: " + action);
               }
            }
         }
      } catch (Throwable e) {
         ward.setCause(e);
      }

      try {
         if (ward.isNull() && asserter != null) {
            asserter.verify();
         }
      } catch (Throwable e) {
         ward.setCause(e);
      }

      ward.escalateCause();
   }

   /**
    * NOTE: These methods should not be exposed -
    * it is breaking SUITA rules.
    * It is due to need introduced by the re-use of
    * VC UI QE lib.
    */
   public FlashSelenium getFlashSelenium() {
      return seleLinks.flashSelenium();
   }

   /**
    * NOTE: These methods should not be exposed -
    * it is breaking SUITA rules.
    * It is introduced to be able to handle login
    * NGC page as it is not flex app.
    * It shouldn't be used on another places outside S?UITA
    * except LoginView and will be fixed once H5 support is
    * introduced.
    */
   public Selenium getSelenium() {
      return seleLinks.selenium();
   }

   /**
    * NOTE: These methods should not be exposed -
    * it is breaking SUITA rules.
    * It is introduced to be able to handle login
    * NGC page as it is not flex app.
    * It shouldn't be used on another places outside S?UITA
    * except LoginView and will be fixed once H5 support is
    * introduced.
    */
   public WebDriver getWebDriver() {
      return seleLinks.driver();
   }

   private DirectID getDirectId(ComponentMatcher baseComponent) {
      return SeleComponentID.toDirectId(baseComponent);
   }
}
