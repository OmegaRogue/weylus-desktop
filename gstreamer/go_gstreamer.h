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

#pragma once

#include <gst/gst.h>
#include <gtk/gtk.h>


//GstBus *bus;
//GstElement *pipeline;

//GdkPaintable *initialize();
//static void shutdown();
void gstreamer_init();
bool gstreamer_bin_add(GstElement *bin, GstElement *elem);
GstElement *gstreamer_pipeline_new (const char *name);
int gstreamer_element_set_state(GstElement *element, int state);