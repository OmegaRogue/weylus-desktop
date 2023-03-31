// Copyright © 2023 omegarogue
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

//
// Created by omegarogue on 31.03.23.
//

#include "go_gstreamer.h"
#include <stdio.h>

#include <stdlib.h>
#include <glib-object.h>
#include <gtk/gtk.h>


//static void shutdown (){
//    /* Free resources */
//    gst_object_unref (bus);
//    gst_element_set_state (pipeline, GST_STATE_NULL);
//    gst_object_unref (pipeline);
//
//}


void gstreamer_init() {
    gst_init (NULL, NULL);
}

bool gstreamer_bin_add(GstElement *bin, GstElement *elem) {
    return gst_bin_add(GST_BIN(bin), elem);
}
int gstreamer_element_set_state(GstElement *element, int state) {
    return gst_element_set_state(element, state);
}
GstElement *gstreamer_pipeline_new (const char *name) {
    return gst_pipeline_new(name);
}

//
//GdkPaintable *initialize() {
//    GstElement *source, *convert, *sink;
//    GdkPaintable *paintable;
//
//        GstStateChangeReturn ret;
//
//
//        /* Initialize GStreamer */
//        gst_init (NULL, NULL);
//
//        /* Create the elements */
//        source = gst_element_factory_make ("videotestsrc", "source");
//        convert = gst_element_factory_make ("videoconvert", "convert");
//        sink = gst_element_factory_make ("gtk4paintablesink", "sink");
//
//
//        g_object_get(sink, "paintable", &paintable, NULL);
//
//        /* Create the empty pipeline */
//    GstElement *pipeline = gst_pipeline_new("test-pipeline");
//        if (!pipeline || !source || !sink) {
//            g_printerr ("Not all elements could be created.\n");
//            exit(-1);
//        }
//        /* Build the pipeline */
//        gst_bin_add_many (GST_BIN (pipeline), source, convert, sink, NULL);
//        if (gst_element_link_many(source, convert, sink, NULL) != TRUE) {
//            g_printerr ("Elements could not be linked.\n");
//            gst_object_unref (pipeline);
//            exit(-1);
//        }
//        /* Modify the source's properties */
//        g_object_set (source, "pattern", 0, NULL);
//
//        /* Start playing */
//        ret = gst_element_set_state (pipeline, GST_STATE_PLAYING);
//        if (ret == GST_STATE_CHANGE_FAILURE) {
//            g_printerr ("Unable to set the pipeline to the playing state.\n");
//            gst_object_unref (pipeline);
//            exit(-1);
//        }
//}