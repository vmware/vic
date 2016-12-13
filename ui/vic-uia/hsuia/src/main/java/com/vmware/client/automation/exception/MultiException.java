/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.exception;

import java.io.PrintStream;
import java.io.PrintWriter;
import java.util.LinkedList;
import java.util.List;

/**
 * Wraps multiple exceptions. Allows multiple exceptions to be thrown as a
 * single exception.
 *
 */
@SuppressWarnings("serial")
public class MultiException extends Exception {
   private List<Throwable> wrappedThrowables = new LinkedList<Throwable>();

   public MultiException() {
      super("Multiple exceptions");
   }

   public MultiException(String message) {
      super(message);
   }

   /**
    * Adds another exception or adds several exceptions, part of a
    * MultiException argument.
    *
    * @param t
    *           Throwable instance to add. Can be also MultiException instance.
    */
   public void add(Throwable t) {
      if (t instanceof MultiException) {
         MultiException me = (MultiException) t;
         for (Throwable e : me.getWrappedThrowables()) {
            wrappedThrowables.add(e);
         }
      } else
         wrappedThrowables.add(t);
   }

   /**
    * Returns a list of the wrapped exceptions.
    *
    * @return List<Throwable>
    */
   public List<Throwable> getWrappedThrowables() {
      List<Throwable> list = wrappedThrowables;
      return list;
   }

   /**
    * Throw a MultiException. If this multi exception is empty then no action is
    * taken. If it contains a single exception that is thrown, otherwise the
    * multi exception is thrown.
    *
    * @exception Exception
    */
   public void ifExceptionThrow() throws Exception {
      switch (wrappedThrowables.size()) {
         case 0:
            break;
         case 1:
            Throwable th = wrappedThrowables.get(0);
            if (th instanceof Error)
               throw (Error) th;
            if (th instanceof Exception)
               throw (Exception) th;
         default:
            // Throw MultiException with modified message and stack trace
            // containing details for the wrapped exceptions
            MultiException me = new MultiException(this.toString());
            me.add(this);
            me.setStackTrace(wrappedThrowables.get(0).getStackTrace());
            throw me;
      }
   }

   /**
    * Returns a short description of this MultiException, containing the
    * descriptions of its wrapped exceptions.
    */
   @Override
   public String toString() {
      if (wrappedThrowables.size() > 0)
         return MultiException.class.getSimpleName() + wrappedThrowables;
      return MultiException.class.getSimpleName() + "[]";
   }

   /**
    * Prints this MultiException and its wrapped exceptions, along with their
    * stack traces, to the standard error stream.
    */
   @Override
   public void printStackTrace() {
      printStackTrace(System.err);
   }

   /**
    * Prints this MultiException and its wrapped exceptions, along with their
    * stack traces, to the specified print stream.
    */
   @Override
   public void printStackTrace(PrintStream out) {
      out.println(this);
      for (Throwable t : wrappedThrowables)
         t.printStackTrace(out);
   }

   /**
    * Prints this MultiException and its wrapped exceptions, along with their
    * stack traces, to the specified print writer.
    */
   @Override
   public void printStackTrace(PrintWriter out) {
      out.println(this);
      for (Throwable t : wrappedThrowables)
         t.printStackTrace(out);
   }
}
