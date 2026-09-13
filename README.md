<div align="center">
  <img src="almost-logo.png" height="120">
  <h1 align="center">almost</h1>
  <p align="center">An on-demand immutability based on file attributes.</p>
</div>

# almost
An on-demand immutability for VanillaOS.

> **Note**: This is a work in progress. It is not ready for production use.

### Read here
This program is meant to be used with [apx](https://github.com/vanilla-os/apx), 
an apt replacement for VanillaOS.

### Help
```
Usage: 
almost [options] [command]

Options:
	--help/-h		show this message
	--verbose/-v		show more verbosity
	--version/-V		show version

Commands:
	enter			set the filesystem as ro or rw until reboot
	config			show the current configuration
	check			check whether the filesystem is read-only or read-write
```

## Use of Generative AI

Maintainers may use generative AI tools as assistants while working on almost. Non-trivial assisted commits disclose the tool, model, and scope of the work.

AI tools may assist with code comments, documentation, repetitive code, and issue triage. Maintainers make project decisions and review every assisted change before it is merged.

Use these trailers for non-trivial assisted commits:

```plain
Assisted-by: <tool>:<model-version>
AI-Scope: <what the tool generated and the prompt or a short prompt summary>
```

Single-line completions, renames, and formatting changes do not need trailers.

Coding agents must also follow [AGENTS.md](AGENTS.md) before changing files,
creating commits, or opening pull requests.
