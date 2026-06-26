//! The ANSI-default theme.
//!
//! [`Token`]s are semantic slots. A [`Theme`] resolves each token to a **named ANSI
//! color** taken from the terminal's 16-color palette — never a raw RGB value by
//! default. The terminal (and its theme manager, e.g. tinty) therefore owns the actual
//! rendering, and every byteowlz tool automatically matches the operator's light/dark
//! theme and palette with zero per-tool configuration.
//!
//! Code references tokens, never raw colors. This is the consistency layer for the
//! whole family, and the reason yazi/helix/lazygit look right in any terminal.

use ratatui::style::{Color, Modifier, Style};

/// A semantic color slot. Code references these instead of raw colors.
///
/// Each variant maps to one role; see [`Theme::color`] for the resolution.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum Token {
    /// Background; resolves to the terminal's default background (`Reset`).
    Surface,
    /// Primary text; resolves to the terminal's default foreground (`Reset`).
    Primary,
    /// Secondary text, metadata, hints.
    Muted,
    /// The single saturated color reserved for focus and primary actions.
    Accent,
    /// Positive state (ok, success).
    Success,
    /// Destructive state (delete, error).
    Danger,
    /// Cautionary state.
    Warning,
    /// Neutral informational state.
    Info,
}

/// A theme: a mapping from [`Token`] to color, applied identically across every
/// screen and every tool.
///
/// Use [`Theme::ansi_default`] unless you have a stated reason to override (a
/// deliberately branded tool). Designing with raw hex is the smell; pinning a token to
/// RGB is an opt-in exception you state explicitly via [`Theme::with_token`].
#[derive(Debug, Clone, Copy)]
pub struct Theme {
    surface: Color,
    primary: Color,
    muted: Color,
    accent: Color,
    success: Color,
    danger: Color,
    warning: Color,
    info: Color,
}

impl Theme {
    /// The default ANSI theme. Defers every color to the terminal palette.
    ///
    /// `Surface`/`Primary` use `Reset` (terminal default bg/fg); the rest use the
    /// terminal's named palette colors. There is exactly one `Accent`.
    #[must_use]
    pub const fn ansi_default() -> Self {
        Self {
            surface: Color::Reset,
            primary: Color::Reset,
            muted: Color::DarkGray,
            accent: Color::Blue,
            success: Color::Green,
            danger: Color::Red,
            warning: Color::Yellow,
            info: Color::Cyan,
        }
    }

    /// Override a single token's color.
    ///
    /// Use sparingly and with a stated reason. Pinning `Accent` to an RGB for a
    /// branded tool is the canonical legitimate use.
    #[must_use]
    pub const fn with_token(mut self, token: Token, color: Color) -> Self {
        match token {
            Token::Surface => self.surface = color,
            Token::Primary => self.primary = color,
            Token::Muted => self.muted = color,
            Token::Accent => self.accent = color,
            Token::Success => self.success = color,
            Token::Danger => self.danger = color,
            Token::Warning => self.warning = color,
            Token::Info => self.info = color,
        }
        self
    }

    /// Resolve a token to its color.
    #[must_use]
    pub const fn color(self, token: Token) -> Color {
        match token {
            Token::Surface => self.surface,
            Token::Primary => self.primary,
            Token::Muted => self.muted,
            Token::Accent => self.accent,
            Token::Success => self.success,
            Token::Danger => self.danger,
            Token::Warning => self.warning,
            Token::Info => self.info,
        }
    }

    /// A plain foreground style for the token (no modifiers).
    #[must_use]
    pub const fn fg(self, token: Token) -> Style {
        Style::new().fg(self.color(token))
    }

    /// A bold foreground style for the token — the primary hierarchy mechanism.
    #[must_use]
    pub const fn fg_bold(self, token: Token) -> Style {
        self.fg(token).add_modifier(Modifier::BOLD)
    }

    /// A dimmed style — for secondary/metadata text.
    #[must_use]
    pub const fn dim(self) -> Style {
        Style::new().add_modifier(Modifier::DIM)
    }

    /// An accented, reversed style — for the focused/active row of a list or menu.
    #[must_use]
    pub const fn focus(self) -> Style {
        self.fg(Token::Accent).add_modifier(Modifier::BOLD)
    }
}

impl Default for Theme {
    fn default() -> Self {
        Self::ansi_default()
    }
}
