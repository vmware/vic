package com.vmware.suitaf.util;

import java.lang.reflect.Array;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * This wrapper class is used in methods that apply the "Typed
 * Parameterization Scheme". In this scheme the parameter values are received
 * through a varargs parameter of type {@link Object}. Then all values are
 * grouped and assigned to local variables according to their intrinsic
 * type. All typed parameters of one method that receive values through the
 * varargs input parameter must be of different type.
 * <br><br>
 * The class constructor transforms the input array into list. It removes
 * the <b>null</b> values. If an instance in the parameter array appears
 * to be an array of {@link Object}s - it is considered a sublist of
 * parameter values. That is why it is split on separate values and they
 * are entered one by one in the parameter list.
 * In this case the parameters resulting from the split are put in
 * the same order in the parameter list.
 * <br>
 * @author dkozhuharov
 */
public class TypedParams {

    /**
     * This interface helps methods with typed parameterization to have
     * postponed calculation of input parameters.<br>
     * This is very helpful in cases when some of the parameters are very
     * resource consuming to produce, but are needed only in some of the
     * method's invocations (e.g. assertion methods).<br>
     * <br>
     * @author dkozhuharov
     */
    public interface ParamsProvider {
        public Object[] getParams();
    }

    private final List<Object> paramList;

    /**
     * This constructor transforms the input array into list and assigns
     * it to a private final list field. Transformation is performed by
     * method {@link TypedParams#toList(Object[])}.
     * @param paramArray - the input array that contains
     * @param invokeProviders - determines if {@link ParamsProvider}
     * instances contained in the array should be invoked to produce their
     * parameter lists. This parameter will be <b>true</b> if the caller
     * is the final consumer of the parameter values and is about to
     * use the incoming parameters. This parameter will be <b>false</b> if
     * the caller is a "transit" method that will just alter the typed
     * parameters list and pass the execution further.
     *
     */
    public TypedParams(boolean invokeProviders, Object...paramArray) {
        paramList = toList(invokeProviders, paramArray);
    }

    /**
     * The method transforms the input array into list. It removes
     * the <b>null</b> values. If an instance in the parameter array appears
     * to be an array of {@link Object}s - it is considered a sublist of
     * parameter values. That is why it is split on separate values and they
     * are entered one by one in the parameter list.
     * In this case the parameters resulting from the split are put in
     * the same order in the parameter list.
     *
     * @param paramArray - an array containing parameters or arrays of
     * parameters.
     * @return - a normalized list of parameters
     */
    public static List<Object> toList(
            boolean invokeProviders, Object...paramArray) {
        List<Object> paramList = new ArrayList<Object>();

        // Safe behavior on null input
        if (paramArray != null) {
            paramList.addAll(Arrays.asList(paramArray));
        }

        flattenList(invokeProviders, paramList);

        return paramList;
    }
    public static Object[] toArray(
            boolean invokeProviders, Object...paramArray) {
        return toList(invokeProviders, paramArray).toArray();
    }

    private static void flattenList(
            boolean invokeProviders, List<Object> paramList) {
        int i=0;
        while (i<paramList.size()) {
            Object par = paramList.get(i);

            // Skip nulls
            if (par == null) {
                paramList.remove(i);
                continue;
            }

            // If parameter providers invocation is allowed and a
            // parameter provider is found - get the provider parameters
            // and push them back in the stack keeping their order
            if (invokeProviders && (par instanceof ParamsProvider)) {
                paramList.remove(i);
                paramList.add(
                        ((ParamsProvider) par).getParams()
                );
                continue;
            }

            // If an array parameter is found - break it up on elements
            // and push them back in the stack keeping their order
            if (par.getClass().isArray()) {
                paramList.remove(i);
                for (int j=(Array.getLength(par)-1); j>=0; j--) {
                    paramList.add(i, Array.get(par, j));
                }
                continue;
            }

            i++;
        }
    }

    /** Returns the current state of the wrapped list of parameters
     * @return list of typed parameters
     */
    public List<Object> toList() {
        return paramList;
    }

    /** Returns the current state of the wrapped list of parameters as array.
     * @return array of typed parameters
     */
    public Object[] toArray() {
        return paramList.toArray();
    }

    /**
     * This is a convenience version of method
     * {@link TypedParams#extractFirst(Class, Object)}
     * The method removes the parameters of requested class from the wrapped
     * parameter list and returns the first of them. If no parameter of
     * the requested class was found returns <b>null</b>.
     *
     * @param <B> class of the parameter to extract from the parameter list
     * @param paramClass - instance of the class of parameters
     * to extract from the wrapped parameter list.
     * @return the first parameter of the requested class
     */
    public <B> B extractFirst(Class<B> paramClass) {
        return extractFirst(paramClass, null);
    }
    /**
     * The method removes the parameters of requested type from the wrapped
     * parameter list and returns the first of them. If no parameter of
     * the requested type was found returns the default value provided.
     *
     * The method removes the parameters of requested class from the wrapped
     * parameter list and returns the first of them. If no parameter of
     * the requested class was found returns <b>null</b>.
     *
     * @param <B> class of the parameter to extract from the parameter list
     * @param paramClass - instance of the class of parameters
     * to extract from the wrapped parameter list.
     * @param defaultValue - the value to be returned if no parameter
     * of the requested class was found
     * @return the first parameter of the requested class
     */
    public <B> B extractFirst(Class<B> paramClass, B defaultValue) {
        List<B> params = extractAll(paramClass);
        if (params.size() == 0)
            return null;
        else
            return params.get(0);
    }
    /**
     * This is a convenience version of method
     * {@link TypedParams#extractLast(Class, Object)}
     * The method removes the parameters of requested class from the wrapped
     * parameter list and returns the last of them. If no parameter of
     * the requested class was found returns <b>null</b>.
     *
     * @param <B> class of the parameter to extract from the parameter list
     * @param paramClass - instance of the class of parameters
     * to extract from the wrapped parameter list.
     * @return the last parameter of the requested class
     */
    public <B> B extractLast(Class<B> paramClass) {
        return extractLast(paramClass, null);
    }
    /**
     * The method removes the parameters of requested class from the wrapped
     * parameter list and returns the last of them. If no parameter of
     * the requested class was found returns <b>null</b>.
     *
     * @param <B> class of the parameter to extract from the parameter list
     * @param paramClass - instance of the class of parameters
     * to extract from the wrapped parameter list.
     * @param defaultValue - the value to be returned if no parameter
     * of the requested class was found
     * @return the last parameter of the requested class
     */
    public <B> B extractLast(Class<B> paramClass, B defaultValue) {
        List<B> params = extractAll(paramClass);
        if (params.size() == 0)
            return defaultValue;
        else
            return params.get(params.size()-1);
    }
    /**
     * The method removes the parameters of requested class from the wrapped
     * parameter list and returns a list of them. If no parameter of
     * the requested class was found returns empty list.
     *
     * @param <B> class of the parameter to extract from the parameter list
     * @param paramClass - instance of the class of parameters
     * to extract from the wrapped parameter list.
     * @return a list of parameters of the requested class
     */
    @SuppressWarnings("unchecked")
    public <B> List<B> extractAll(Class<B> paramClass) {
        ArrayList<B> paramsOfType = new ArrayList<B>();

        int i = 0;
        while (i < paramList.size()) {
            if (paramClass.isInstance(paramList.get(i))) {
                paramsOfType.add((B) paramList.get(i));
                paramList.remove(i);
            }
            else {
                i++;
            }
        }

        return paramsOfType;
    }

    /**
     * Appends a list of typed parameters to the end of the wrapped
     * parameter list.
     *
     * @param params - list of typed parameters to append
     */
    public void append(List<?> params) {
        paramList.addAll(params);
    }

    /**
     * Appends a typed parameter to the end of the wrapped parameter list.
     *
     * @param params - list of typed parameters to append
     */
    public void append(Object param) {
        paramList.add(param);
    }

    public void assertAllUsed() {
        if (paramList.size() > 0) {
            StringBuilder sb = new StringBuilder();

            for (Object obj: paramList) {
                sb.append(obj.getClass().getSimpleName());
                sb.append("{" + obj + "} ");
            }

            if (paramList.size() > 0) {
                throw new RuntimeException(
                        "Failed to use following <TypedParams>: " + sb);
            }
        }
    }
}