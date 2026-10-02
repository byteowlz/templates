//! Activate the design-system look inside gpui-kit's theme system.
//!
//! The `{{project_name}}-design` crate emits a gpui-component `ThemeSet` JSON
//! document (schemes -> closed roles -> theme tokens). This module loads it
//! into the component `ThemeRegistry` and applies the configured scheme.

use gpui_kit::component::{Theme, ThemeRegistry};
use gpui_kit::{Global, SharedString};
use {{project_name}}_design::example_theme_set_json;

use crate::config::Config;

/// The example scheme names shipped by the design crate. Replace with your own
/// scheme names once you swap the example look for a real one.
pub const DARK: &str = "Example Dark";
pub const LIGHT: &str = "Example Light";

/// Faces that are not part of gpui-component's theme: the accent face used on
/// section headings and the status bar.
pub struct Fonts {
    pub accent: SharedString,
}

impl Global for Fonts {}

/// The accent face, from config (default: the design identity's).
pub fn accent_font(cx: &gpui_kit::App) -> SharedString {
    cx.try_global::<Fonts>()
        .map(|f| f.accent.clone())
        .unwrap_or_else(|| SharedString::from("SystemUIFont"))
}

/// Native pixel size of the accent face; pixel faces blur at other sizes.
pub const ACCENT_PX: f32 = 11.0;

/// Load the theme set and activate the configured scheme.
pub fn init(cx: &mut gpui_kit::App, config: &Config) {
    cx.set_global(Fonts {
        accent: SharedString::from(config.accent_font.clone()),
    });
    match example_theme_set_json() {
        Ok(json) => {
            if let Err(e) = ThemeRegistry::global_mut(cx).load_themes_from_str(&json) {
                log::error!("loading theme set: {e:#}");
            }
        }
        Err(e) => log::error!("building theme set: {e:#}"),
    }
    apply(cx, &config.theme);
}

/// Activate a registered theme by name.
///
/// gpui-component keeps two theme globals: its own `Theme`, and a copy in
/// gpui-base that paints markdown, code blocks, scrollbars and resize handles.
/// `apply_config` only updates the former; `Theme::change` is what refreshes
/// both, so the chosen scheme is installed as the mode's theme and then
/// activated through `change`.
pub fn apply(cx: &mut gpui_kit::App, name: &str) {
    let Some(config) = ThemeRegistry::global_mut(cx).themes().get(name).cloned() else {
        log::error!("theme {name:?} is not registered");
        return;
    };
    let mode = config.mode;
    {
        let theme = Theme::global_mut(cx);
        if mode.is_dark() {
            theme.dark_theme = config;
        } else {
            theme.light_theme = config;
        }
    }
    Theme::change(mode, None, cx);
    cx.refresh_windows();
}

/// Toggle dark/light and return the activated name.
pub fn toggle(cx: &mut gpui_kit::App, current: &str) -> &'static str {
    let next: &'static str = if current == DARK { LIGHT } else { DARK };
    apply(cx, next);
    next
}

#[cfg(test)]
mod tests {
    use gpui_kit::component::ThemeSet;

    /// The emitted theme set must parse with gpui-component's own schema and
    /// carry role-mapped colors under the keys the toolkit reads.
    use super::*;
    #[test]
    fn emitted_theme_set_parses_with_toolkit_schema() {
        let json = example_theme_set_json().unwrap();
        let set: ThemeSet = serde_json::from_str(&json).unwrap();
        assert_eq!(set.themes.len(), 2);
        let dark = set.themes.iter().find(|t| t.name == DARK).unwrap();
        assert_eq!(dark.mode.name(), "dark");
        assert!(dark.colors.background.is_some());
        assert!(dark.colors.primary.is_some());
        assert!(dark.colors.sidebar.is_some());
        assert!(dark.highlight.is_some(), "syntax theme must land");
    }
}