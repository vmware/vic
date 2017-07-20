// ------------------------------------------------------------------------------
// Javascript initialization required for using the HTML Client SDK API.
// This also creates the plugin's private namespace com_vmware_vic
//
//                       DO NOT MAKE CHANGES HERE!
// ------------------------------------------------------------------------------

// WEB_PLATFORM is the VMware Web Client platform reference:
// - When the plugin runs in the HTML client WEB_PLATFORM is a real JS object
// - When it runs in the Flex client WEB_PLATFORM is defined as the Flash container

var WEB_PLATFORM = self.parent.WEB_PLATFORM;
var isChromeBrowser = (window.navigator.userAgent.indexOf("Chrome/") >= 0);
var isFlexClient = !!self.parent.document.getElementById("container_app");

if (!WEB_PLATFORM || (isChromeBrowser && isFlexClient)) {
      WEB_PLATFORM = self.parent.document.getElementById("container_app");
      if (isChromeBrowser) {
            // Object.create is required to support Chrome version >= 55 on Flex client
            WEB_PLATFORM = Object.create(WEB_PLATFORM);
      }
      if (WEB_PLATFORM === null) {
            WEB_PLATFORM = {};
      }
      self.parent.WEB_PLATFORM = WEB_PLATFORM;

      // The web context starts with a different root path depending on which client is running.
      if (!WEB_PLATFORM.getRootPath) {
            WEB_PLATFORM.getRootPath = function () {
                  return "/vsphere-client";
            }
      }
      // Declare unknown client type explicitly.
      if (!WEB_PLATFORM.getClientType) {
            WEB_PLATFORM.getClientType = function () {
                  return "flex";
            }
      }
      // Declare unknown client version explicitly.
      if (!WEB_PLATFORM.getClientVersion) {
            WEB_PLATFORM.getClientVersion = function () {
                  return "6.0";
            }
      }
}

// Define a private namespace using the plugin bundle name,
// It should be the only global symbol added by this plugin!
var com_vmware_vic;
if (!com_vmware_vic) {
      com_vmware_vic = {};

      // The web context path to use for server requests, compatible with Flex and HTML clients
      com_vmware_vic.webContextPath = WEB_PLATFORM.getRootPath() + "/vic";

      // The API setup is done inside an anonymous function to keep things clean.
      // See the SDK documentation for more info on those APIs.
      (function () {
            // Namespace shortcut
            var ns = com_vmware_vic;

            // ------------------------ Private functions -------------------------------

            // Get a string from the resource bundle defined in plugin.xml
            function getString(key, params) {
                  return WEB_PLATFORM.getString("com_vmware_vic", key, params);
            }

            // Get a parameter value from the current document URL
            function getURLParameter(name) {
                  return (new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)')
                        .exec(location.href) || [, ""])[1].replace(/\+/g, '%20') || null;
            }

            // Build the REST url prefix to retrieve a list of properties,
            // this is mapped to the DataAccessController on the java side.
            function buildDataUrl(objectId, propList) {
                  var propStr = propList.toString();
                  var dataUrl = ns.webContextPath +
                        "/rest/data/properties/" + objectId + "?properties=" + propStr;
                  return dataUrl;
            }

            // -------------------------- Public APIs --------------------------------

            // Functions exported to the com_vmware_vic namespace
            ns.getString = getString;
            ns.buildDataUrl = buildDataUrl;

            // Get the current context object id or return null if none is defined
            WEB_PLATFORM.getObjectId = function () {
                  return getURLParameter("objectId");
            };
            // Get the current action Uid or return null if none is defined
            WEB_PLATFORM.getActionUid = function () {
                  return getURLParameter("actionUid");
            };
            // Get the comma-separated list of object ids for an action, or null for a global action
            WEB_PLATFORM.getActionTargets = function () {
                  return getURLParameter("targets");
            };
            // Get the current locale
            WEB_PLATFORM.getLocale = function () {
                  return getURLParameter("locale");
            };

            // Get the info provided in a global view using a vCenter selector
            WEB_PLATFORM.getVcSelectorInfo = function () {
                  var info = {
                        serviceGuid: getURLParameter("serviceGuid"),
                        sessionId: getURLParameter("sessionId"),
                        serviceUrl: getURLParameter("serviceUrl")
                  };
                  return info;
            };

            // Set a refresh handler called when the user hits Refresh in the WebClient top toolbar.
            // This is the Flex Client implementation which doesn't need document as 2nd argument.
            if (WEB_PLATFORM.getClientType() === "flex") {
                  WEB_PLATFORM.setGlobalRefreshHandler = function (handler) {
                        WEB_PLATFORM["refresh" + window.name] = handler;
                  };
            }
      })();
} // end of if (!com_vmware_vic)
