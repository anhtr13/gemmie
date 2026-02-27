# Gemmie

🎶 Gimme Gimme Gimme your love after midnight 🎶.
Just kidding, it's a tool to bring GeminiAi to the terminal.

## Installation

**1. Install the repository:**

```bash
go install github.com/anhtr13/gemmie@latest
```

**2. Setup Gemini API Key:**

- Get API key:
  - Go to [Google AI Studio](https://aistudio.google.com)
  - Sign in with Google account
  - "Create API Key"
  - Copy the generated key

- Save API key:

  ```sh
  gemmie apikey [api_key]
  ```

## Usage

`gemmie --help` for the manual

**Example:**

```bash
gemmie --model=gemini-3.0-pro --lang=Vietnamese --temp=2.0 --limit=6900 -p="write a story about a magic backpack."
```

_**Note:** Never share your API key._
