//! The root application view. Shows the design-system role vocabulary in a
//! real, interactive layout: a sunken sidebar, an elevated card, primary and
//! outline buttons, muted/secondary text, status badges and a dark/light
//! toggle — all read from `cx.theme()` (the role layer), never a
//! host-specific hex.

use gpui_kit::component::button::{Button, ButtonVariants as _};
use gpui_kit::component::ActiveTheme as _;
use gpui_kit::{
    ClickEvent, Context, FontWeight, Hsla, IntoElement, ParentElement as _, Render, Styled as _,
    Window, div, px, transparent_black,
};

use crate::config::Config;
use crate::theme;

pub struct App {
    theme_name: String,
    count: u32,
}

impl App {
    pub fn new(config: Config, _window: &mut Window, _cx: &mut Context<Self>) -> Self {
        Self {
            theme_name: config.theme.clone(),
            count: 0,
        }
    }

    fn toggle_theme(&mut self, cx: &mut Context<Self>) {
        self.theme_name = theme::toggle(cx, &self.theme_name).to_string();
        cx.notify();
    }
}

impl Render for App {
    fn render(&mut self, _window: &mut Window, cx: &mut Context<Self>) -> impl IntoElement {
        let theme = cx.theme().clone();
        let accent_font = theme::accent_font(cx);
        let is_dark = self.theme_name == theme::DARK;

        div()
            .w_full()
            .h_full()
            .flex()
            .flex_col()
            .bg(theme.background)
            .text_color(theme.foreground)
            // Title bar sits on the sunken surface (surface-sunken role).
            .child(
                div()
                    .flex()
                    .items_center()
                    .gap_3()
                    .px_4()
                    .h(px(44.))
                    .bg(theme.sidebar)
                    .border_b_1()
                    .border_color(theme.border)
                    .child(
                        div()
                            .text_size(px(theme::ACCENT_PX))
                            .font_family(accent_font.clone())
                            .text_color(theme.primary)
                            .child("{{project_name}}"),
                    )
                    .child(div().flex_1())
                    .child(
                        Button::new("toggle-theme")
                            .outline()
                            .label(if is_dark { "Light mode" } else { "Dark mode" })
                            .on_click(cx.listener(|this, _: &ClickEvent, _window, cx| {
                                this.toggle_theme(cx)
                            })),
                    ),
            )
            // Body: sunken sidebar + elevated content card.
            .child(
                div()
                    .flex_1()
                    .flex()
                    .min_h_0()
                    .child(
                        div()
                            .w(px(220.))
                            .bg(theme.sidebar)
                            .border_r_1()
                            .border_color(theme.border)
                            .p_3()
                            .flex()
                            .flex_col()
                            .gap_2()
                            .child(sidebar_item(theme.muted_foreground, "Home", true))
                            .child(sidebar_item(theme.muted_foreground, "Activity", false))
                            .child(sidebar_item(theme.muted_foreground, "Settings", false)),
                    )
                    .child(
                        div()
                            .flex_1()
                            .p_6()
                            .flex()
                            .flex_col()
                            .gap_4()
                            .min_w_0()
                            .child(
                                div()
                                    .text_size(px(22.))
                                    .font_weight(FontWeight::BOLD)
                                    .font_family(accent_font.clone())
                                    .child("Design-system roles"),
                            )
                            .child(
                                div()
                                    .text_color(theme.muted_foreground)
                                    .child(
                                        "Every color below comes from the scheme's closed role \
                                         layer — reskin the scheme and the whole app follows. \
                                         No host-specific hex.",
                                    ),
                            )
                            // A raised card uses the `surface` role.
                            .child(
                                div()
                                    .bg(theme.popover)
                                    .border_1()
                                    .border_color(theme.border)
                                    .rounded(px(theme.radius_lg.as_f32()))
                                    .p_4()
                                    .flex()
                                    .flex_col()
                                    .gap_3()
                                    .child(
                                        div()
                                            .text_size(px(14.))
                                            .font_weight(FontWeight::SEMIBOLD)
                                            .child("An elevated card"),
                                    )
                                    .child(
                                        div()
                                            .flex()
                                            .items_center()
                                            .gap_3()
                                            .child(
                                                Button::new("greet")
                                                    .primary()
                                                    .label("Say hi")
                                                    .on_click(cx.listener(|this, _: &ClickEvent, _window, cx| {
                                                        this.count += 1;
                                                        cx.notify();
                                                    })),
                                            ),
                                    )
                                    .child(
                                        div()
                                            .text_color(theme.muted_foreground)
                                            .child(format!("Greeted {} time(s).", self.count)),
                                    ),
                            )
                            .child(
                                div()
                                    .flex()
                                    .gap_2()
                                    .child(status_badge(theme.success, "ok"))
                                    .child(status_badge(theme.info, "info"))
                                    .child(status_badge(theme.warning, "warn"))
                                    .child(status_badge(theme.danger, "error")),
                            ),
                    ),
            )
            .child(
                div()
                    .h(px(26.))
                    .px_4()
                    .flex()
                    .items_center()
                    .bg(theme.sidebar)
                    .border_t_1()
                    .border_color(theme.border)
                    .text_size(px(11.))
                    .text_color(theme.muted_foreground)
                    .child(format!("{} · {} ready", env!("CARGO_PKG_NAME"), self.theme_name)),
            )
    }
}

fn sidebar_item(color: Hsla, label: &str, active: bool) -> impl IntoElement {
    let bg = if active {
        color.alpha(0.12)
    } else {
        transparent_black()
    };
    div()
        .px_2()
        .py_1()
        .rounded(px(4.))
        .bg(bg)
        .text_color(color)
        .child(label.to_string())
}

fn status_badge(color: Hsla, label: &str) -> impl IntoElement {
    div()
        .px_2()
        .py_1()
        .rounded(px(4.))
        .bg(color.alpha(0.16))
        .text_color(color)
        .text_size(px(12.))
        .child(label.to_string())
}