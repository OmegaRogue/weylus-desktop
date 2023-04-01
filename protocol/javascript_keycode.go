/*
 * Copyright © 2023 omegarogue
 * SPDX-License-Identifier: AGPL-3.0-or-later
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package protocol
package protocol

import "github.com/diamondburned/gotk4/pkg/gdk/v4"

type KeyData struct {
	Name     string
	Key      string
	Location KeyboardLocation
}

//goland:noinspection SpellCheckingInspection
var (
	CodeValue = map[uint]string{
		0x0009: "Escape",
		0x000A: "Digit1",
		0x000B: "Digit2",
		0x000C: "Digit3",
		0x000D: "Digit4",
		0x000E: "Digit5",
		0x000F: "Digit6",
		0x0010: "Digit7",
		0x0011: "Digit8",
		0x0012: "Digit9",
		0x0013: "Digit0",
		0x0014: "Minus",
		0x0015: "Equal",
		0x0016: "Backspace",
		0x0017: "Tab",
		0x0018: "KeyQ",
		0x0019: "KeyW",
		0x001A: "KeyE",
		0x001B: "KeyR",
		0x001C: "KeyT",
		0x001D: "KeyY",
		0x001E: "KeyU",
		0x001F: "KeyI",
		0x0020: "KeyO",
		0x0021: "KeyP",
		0x0022: "BracketLeft",
		0x0023: "BracketRight",
		0x0024: "Enter",
		0x0025: "ControlLeft",
		0x0026: "KeyA",
		0x0027: "KeyS",
		0x0028: "KeyD",
		0x0029: "KeyF",
		0x002A: "KeyG",
		0x002B: "KeyH",
		0x002C: "KeyJ",
		0x002D: "KeyK",
		0x002E: "KeyL",
		0x002F: "Semicolon",
		0x0030: "Quote",
		0x0031: "Backquote",
		0x0032: "ShiftLeft",
		0x0033: "Backslash",
		0x0034: "KeyZ",
		0x0035: "KeyX",
		0x0036: "KeyC",
		0x0037: "KeyV",
		0x0038: "KeyB",
		0x0039: "KeyN",
		0x003A: "KeyM",
		0x003B: "Comma",
		0x003C: "Period",
		0x003D: "Slash",
		0x003E: "ShiftRight",
		0x003F: "NumpadMultiply",
		0x0040: "AltLeft",
		0x0041: "Space",
		0x0042: "CapsLock",
		0x0043: "F1",
		0x0044: "F2",
		0x0045: "F3",
		0x0046: "F4",
		0x0047: "F5",
		0x0048: "F6",
		0x0049: "F7",
		0x004A: "F8",
		0x004B: "F9",
		0x004C: "F10",
		0x004D: "NumLock",
		0x004E: "ScrollLock",
		0x004F: "Numpad7",
		0x0050: "Numpad8",
		0x0051: "Numpad9",
		0x0052: "NumpadSubtract",
		0x0053: "Numpad4",
		0x0054: "Numpad5",
		0x0055: "Numpad6",
		0x0056: "NumpadAdd",
		0x0057: "Numpad1",
		0x0058: "Numpad2",
		0x0059: "Numpad3",
		0x005A: "Numpad0",
		0x005B: "NumpadDecimal",
		0x005C: "Unidentified",
		0x005D: "Lang5",
		0x005E: "IntlBackslash",
		0x005F: "F11",
		0x0060: "F12",
		0x0061: "IntlRo",
		0x0062: "Lang3",
		0x0063: "Lang4",
		0x0064: "Convert",
		0x0065: "KanaMode",
		0x0066: "NonConvert",
		0x0067: "Unidentified",
		0x0068: "NumpadEnter",
		0x0069: "ControlRight",
		0x006A: "NumpadDivide",
		0x006B: "PrintScreen",
		0x006C: "AltRight",
		0x006D: "Unidentified",
		0x006E: "Home",
		0x006F: "ArrowUp",
		0x0070: "PageUp",
		0x0071: "ArrowLeft",
		0x0072: "ArrowRight",
		0x0073: "End",
		0x0074: "ArrowDown",
		0x0075: "PageDown",
		0x0076: "Insert",
		0x0077: "Delete",
		0x0078: "Unidentified",
		0x0079: "VolumeMute",
		0x007A: "VolumeDown",
		0x007B: "VolumeUp",
		0x007C: "Power",
		0x007D: "NumpadEqual",
		0x007E: "Unidentified",
		0x007F: "Pause",
		0x0080: "Unidentified",
		0x0081: "NumpadComma",
		0x0082: "Lang1",
		0x0083: "Lang2",
		0x0084: "IntlYen",
		0x0085: "MetaLeft",
		0x0086: "MetaRight",
		0x0087: "ContextMenu",
		0x0088: "BrowserStop",
		0x0089: "Again",
		0x008A: "Props",
		0x008B: "Undo",
		0x008C: "Select",
		0x008D: "Copy",
		0x008E: "Open",
		0x008F: "Paste",
		0x0090: "Find",
		0x0091: "Cut",
		0x0092: "Help",
		0x0093: "Unidentified",
		0x0094: "LaunchApp2",
		0x0095: "Unidentified",
		0x0096: "Sleep",
		0x0097: "WakeUp",
		0x0098: "LaunchApp1",
		0x0099: "Unidentified",
		0x009A: "Unidentified",
		0x009B: "Unidentified",
		0x009C: "Unidentified",
		0x009D: "Unidentified",
		0x009E: "Unidentified",
		0x009F: "Unidentified",
		0x00A0: "Unidentified",
		0x00A1: "Unidentified",
		0x00A2: "Unidentified",
		0x00A3: "LaunchMail",
		0x00A4: "BrowserFavorites",
		0x00A5: "Unidentified",
		0x00A6: "BrowserBack",
		0x00A7: "BrowserForward",
		0x00A8: "Unidentified",
		0x00A9: "Eject",
		0x00AA: "Unidentified",
		0x00AB: "MediaTrackNext",
		0x00AC: "MediaPlayPause",
		0x00AD: "MediaTrackPrevious",
		0x00AE: "MediaStop",
		0x00AF: "Unidentified",
		0x00B0: "Unidentified",
		0x00B1: "Unidentified",
		0x00B2: "Unidentified",
		0x00B3: "MediaSelect",
		0x00B4: "BrowserHome",
		0x00B5: "BrowserRefresh",
		0x00B6: "Unidentified",
		0x00B7: "Unidentified",
		0x00B8: "Unidentified",
		0x00B9: "Unidentified",
		0x00BA: "Unidentified",
		0x00BB: "NumpadParenLeft",
		0x00BC: "NumpadParenRight",
		0x00BD: "Unidentified",
		0x00BE: "Unidentified",
		0x00BF: "F13",
		0x00C0: "F14",
		0x00C1: "F15",
		0x00C2: "F16",
		0x00C3: "F17",
		0x00C4: "F18",
		0x00C5: "F19",
		0x00C6: "F20",
		0x00C7: "F21",
		0x00C8: "F22",
		0x00C9: "F23",
		0x00CA: "F24",
		0x00CB: "Unidentified",
		0x00CC: "Unidentified",
		0x00CD: "Unidentified",
		0x00CE: "Unidentified",
		0x00CF: "Unidentified",
		0x00E0: "Unidentified",
		0x00E1: "BrowserSearch",
	}

	KeyValue = map[uint]string{
		gdk.KEY_Escape:           "Escape",
		gdk.KEY_3270_Enter:       "Enter",
		gdk.KEY_ISO_Enter:        "Enter",
		gdk.KEY_Return:           "Enter",
		gdk.KEY_BackSpace:        "Backspace",
		gdk.KEY_Tab:              "Tab",
		gdk.KEY_Caps_Lock:        "CapsLock",
		gdk.KEY_F1:               "F1",
		gdk.KEY_F2:               "F2",
		gdk.KEY_F3:               "F3",
		gdk.KEY_F4:               "F4",
		gdk.KEY_F5:               "F5",
		gdk.KEY_F6:               "F6",
		gdk.KEY_F7:               "F7",
		gdk.KEY_F8:               "F8",
		gdk.KEY_F9:               "F9",
		gdk.KEY_F10:              "F10",
		gdk.KEY_F11:              "F11",
		gdk.KEY_F12:              "F12",
		gdk.KEY_F13:              "F13",
		gdk.KEY_F14:              "F14",
		gdk.KEY_F15:              "F15",
		gdk.KEY_F16:              "F16",
		gdk.KEY_F17:              "F17",
		gdk.KEY_F18:              "F18",
		gdk.KEY_F19:              "F19",
		gdk.KEY_F20:              "F20",
		gdk.KEY_F21:              "F21",
		gdk.KEY_F22:              "F22",
		gdk.KEY_F23:              "F23",
		gdk.KEY_F24:              "F24",
		gdk.KEY_Num_Lock:         "NumLock",
		gdk.KEY_Scroll_Lock:      "ScrollLock",
		gdk.KEY_KP_Enter:         "Enter",
		gdk.KEY_Katakana:         "KanaMode",
		gdk.KEY_3270_PrintScreen: "PrintScreen",
		gdk.KEY_Home:             "Home",
		gdk.KEY_Up:               "ArrowUp",
		gdk.KEY_Page_Up:          "PageUp",
		gdk.KEY_Left:             "ArrowLeft",
		gdk.KEY_Right:            "ArrowRight",
		gdk.KEY_End:              "End",
		gdk.KEY_Down:             "ArrowDown",
		gdk.KEY_Page_Down:        "PageDown",
		gdk.KEY_Insert:           "Insert",
		gdk.KEY_Delete:           "Delete",
		gdk.KEY_AudioMute:        "VolumeMute",
		gdk.KEY_AudioLowerVolume: "VolumeDown",
		gdk.KEY_AudioRaiseVolume: "VolumeUp",
		gdk.KEY_Pause:            "Pause",
		gdk.KEY_Menu:             "ContextMenu",
		gdk.KEY_Cancel:           "Cancel",
		gdk.KEY_Undo:             "Undo",
		gdk.KEY_Copy:             "Copy",
		gdk.KEY_Open:             "Open",
		gdk.KEY_Paste:            "Paste",
		gdk.KEY_Find:             "Find",
		gdk.KEY_Cut:              "Cut",
		gdk.KEY_Help:             "Help",
		gdk.KEY_Mail:             "LaunchMail",
		gdk.KEY_Sleep:            "Sleep",
		gdk.KEY_Control_L:        "Control",
		gdk.KEY_Control_R:        "Control",
		gdk.KEY_Alt_L:            "Alt",
		gdk.KEY_Alt_R:            "Alt",
		gdk.KEY_Meta_L:           "Meta",
		gdk.KEY_Meta_R:           "Meta",
		gdk.KEY_Super_L:          "Super",
		gdk.KEY_Super_R:          "Super",
		gdk.KEY_Hyper_L:          "Hyper",
		gdk.KEY_Hyper_R:          "Hyper",
		gdk.KEY_Shift_L:          "Shift",
		gdk.KEY_Shift_R:          "Shift",
		gdk.KEY_Mode_switch:      "AltGraph",
		gdk.KEY_ISO_Level3_Shift: "AltGraph",
		gdk.KEY_ISO_Level3_Latch: "AltGraph",
		gdk.KEY_ISO_Level3_Lock:  "AltGraph",
		gdk.KEY_ISO_Level5_Shift: "AltGraph",
		gdk.KEY_ISO_Level5_Latch: "AltGraph",
		gdk.KEY_ISO_Level5_Lock:  "AltGraph",
	}
)
