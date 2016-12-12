package com.vmware.suitaf.util;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;

public class CallStackAnalyser {
    public static class CodeZone {
        public final String name;
        final List<String> includeElements;
        final List<CodeZone> includeZones;
        final List<String> excludeElements;
        final List<CodeZone> excludeZones;

        public CodeZone(String name) {
            this.name = name;
            this.includeElements = new ArrayList<String>();
            this.includeZones = new ArrayList<CodeZone>();
            this.excludeElements = new ArrayList<String>();
            this.excludeZones = new ArrayList<CodeZone>();
        }

        public CodeZone setInclusion(Object...includedZoneParts) {
            if (includedZoneParts == null) {
                throw new NullPointerException(
                        "No nulls accepted for inclusion criterion");
            }

            includeElements.clear();
            includeZones.clear();
            for (Object part : includedZoneParts) {
                if (part instanceof Package) {
                    includeElements.add(((Package) part).getName());
                }
                else if (part instanceof Class<?>) {
                    includeElements.add(((Class<?>)part).getName());
                }
                else if (part instanceof CodeZone) {
                    includeZones.add((CodeZone)part);
                }
            }

            return this;
        }

        public CodeZone setExclusion(Object...excludedZoneParts ) {
            if (excludedZoneParts == null) {
                throw new NullPointerException(
                        "No nulls accepted for exclusion criterion");
            }

            excludeElements.clear();
            excludeZones.clear();
            for (Object part : excludedZoneParts) {
                if (part instanceof Package) {
                    excludeElements.add(((Package) part).getName());
                }
                else if (part instanceof Class<?>) {
                    excludeElements.add(((Class<?>)part).getName());
                }
                else if (part instanceof CodeZone) {
                    excludeZones.add((CodeZone)part);
                }
            }

            return this;
        }

        boolean checkIncluded(String className, HashSet<CodeZone> covered) {
            // Set default value for empty include list -> include every one
            if (includeElements.size() == 0 && includeZones.size() == 0) {
                return true;
            }
            for (String include : includeElements) {
                if (className.startsWith(include)) {
                    return true;
                }
            }
            for (CodeZone zone : includeZones) {
                if (zone.includes(className, covered)) {
                    return true;
                }
            }
            return false;
        }

        boolean checkExcludes(String className, HashSet<CodeZone> covered) {
            // Set default value for empty exclude list -> exclude no one
            if (excludeElements.size() == 0 && excludeZones.size() == 0) {
                return false;
            }
            for (String exclude : excludeElements) {
                if (className.startsWith(exclude)) {
                    return true;
                }
            }
            for (CodeZone zone : excludeZones) {
                if (zone.includes(className, covered)) {
                    return true;
                }
            }
            return false;
        }

        public boolean includes(String className) {
            return includes(className, new HashSet<CodeZone>());
        }
        boolean includes(String className, HashSet<CodeZone> covered) {
            if (className == null) {
                throw new NullPointerException(
                        "Class name must not be null.");
            }
            if (covered.contains(this)) {
                throw new RuntimeException(
                        "Closed loop dependency found for zone:" + this);
            }
            else {
                covered.add(this);
            }

            if (checkIncluded(className, covered)) {
                if (!checkExcludes(className, covered)) {
                    return true;
                }
            }

            return false;
        }

        @Override
        public boolean equals(Object obj) {
            if (obj instanceof CodeZone) {
                CodeZone cz = (CodeZone) obj;
                return
                this.includeElements.equals(cz.includeElements) &&
                this.excludeElements.equals(cz.excludeElements);
            }
            return false;
        }

        @Override
        public String toString() {
            return "ZONE:" + name
            + " IN:" + includeElements
            + " EX:" + excludeElements;
        }
    }

    public static class CallPoint {
        public final StackTraceElement invoker;
        public final StackTraceElement invoked;
        public final List<CodeZone> enteredZones;
        public CallPoint(StackTraceElement invoker, StackTraceElement invoked,
                List<CodeZone> invokedZones) {
            this.invoker = invoker;
            this.invoked = invoked;
            this.enteredZones = Collections.unmodifiableList(
                    new ArrayList<CodeZone>(invokedZones));
        }

        public String getInvokerPoint() {
            if (invoker == null)
                return "<top>";

            String className = invoker.getClassName();
            className = className.substring(className.lastIndexOf(".") + 1);
            return className + "#" + invoker.getMethodName() +
                "(" + invoker.getLineNumber() + ")";
        }

        public String getInvokedAction() {
            if (invoked == null)
                return "<bottom>";

            String className = invoked.getClassName();
            className = className.substring(className.lastIndexOf(".") + 1);
            return className + "#" + invoked.getMethodName();
        }

        public String getEnteredZones() {
            StringBuilder zonesMsg = new StringBuilder();

            for (CodeZone zone : enteredZones) {
                if (zonesMsg.length() > 0)
                    zonesMsg.append("/");
                zonesMsg.append(zone.name);
            }
            if (zonesMsg.length() > 0)
                zonesMsg.insert(0, '[').append(']');

            return zonesMsg.toString();
        }

        @Override
        public String toString() {
            return "[" + getInvokerPoint() + "]-->" +
                    getEnteredZones() + getInvokedAction();
        }

        @Override
        public boolean equals(Object obj) {
            if (obj instanceof CallPoint) {
                CallPoint cp = (CallPoint) obj;
                return
                CommonUtils.smartEqual(invoker, cp.invoker) &&
                CommonUtils.smartEqual(invoked, cp.invoked);
            }

            return false;
        }
    }

    public static class Analysis {
        public final List<CallPoint> callPoints;
        public Analysis(List<CallPoint> callPoints) {
            this.callPoints = Collections.unmodifiableList(
                    new ArrayList<CallPoint>(callPoints));
        }

        public CallPoint getEntryPoint(CodeZone zone) {
            for (CallPoint actionPOC : callPoints) {
                if (actionPOC.enteredZones.contains(zone)) {
                    return actionPOC;
                }
            }
            return null;
        }

        public CallPoint getExitPoint() {
            int last = callPoints.size() - 1;
            if (last >= 0) {
                return callPoints.get(last);
            }
            return null;
        }

        public List<String> getFullList() {
            ArrayList<String> fullList = new ArrayList<String>();
            for (CallPoint actionPOC : callPoints) {
                fullList.add(actionPOC.toString());
            }

            return fullList;
        }
    }

    /**
     * Factory method that generates {@link Analysis} instance processing given
     * array of call-stack elements.
     * @param fullStack - the call-stack
     * @param codeMap - a list of {@link CodeZone}s that represent the zones of
     * interest in the code
     * @return the stack analysis instance; It contains the {@link CallPoint}s
     * of entrance to the code zones of interest. Also contains the exit
     * {@link CallPoint} where execution had continued outside the zones.
     */
    public static Analysis process(
            StackTraceElement[] fullStack, CodeZone...codeMap) {

        // If no call-stack was provided or no code zones to analyze
        if (fullStack == null || codeMap.length == 0) {
            return null;
        }

        List<CodeZone> oldZones = new ArrayList<CodeZone>();
        StackTraceElement invoker = null;

        List<CodeZone> newZones = new ArrayList<CodeZone>();
        StackTraceElement invoked = null;

        List<CallPoint> callPoints = new ArrayList<CallPoint>();

        for (int i=fullStack.length-1; i>=-1; i--) {

            if (i > -1) {
                // Check for criteria for skipping of call-stack line
                if (fullStack[i].getMethodName().startsWith("access$")) {
                    continue;
                }

                // find the code zones of the new call-stack element
                invoked = fullStack[i];

                for (CodeZone zone : codeMap) {
                    if (zone.includes(invoked.getClassName())) {
                        newZones.add(zone);
                    }
                }
            }

            // Detection of the exit call point
            if (oldZones.size() > 0 && newZones.size() == 0) {
                // add the exit call point to the call-stack analysis record
                callPoints.add(new CallPoint(invoker, invoked, newZones));

                // Prepare for the next analysis loop
                oldZones.clear();
            }

            // Process old zones and continuing zones
            int old_i = 0;
            while (old_i < oldZones.size()) {
                int new_i = newZones.indexOf(oldZones.get(old_i));
                if (new_i >= 0) {
                    // CallStack is still in this zone
                    newZones.remove(new_i);
                    old_i++;
                }
                else {
                    // Call Stack is no more in this zone
                    oldZones.remove(old_i);
                }
            }

            // Process newly entered zones
            if (newZones.size() > 0) {
                // If new code zone is entered this is substantial call point
                // --> add it to the call-stack analysis record
                callPoints.add(new CallPoint(invoker, invoked, newZones));

                // Prepare for the next analysis loop
                oldZones.addAll(newZones);
                newZones.clear();
            }

            invoker = invoked;
            invoked = null;
        }

        return new Analysis(callPoints);
    }
}
