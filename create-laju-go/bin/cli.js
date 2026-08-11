#!/usr/bin/env node

const { program } = require("commander");
const prompts = require("prompts");
const degit = require("degit");
const path = require("path");
const fs = require("fs");
const { execSync } = require("child_process");

/**
 * Files/dirs to strip from the scaffolded project.
 * These are repo-internal / dev-only artifacts that don't belong in a fresh scaffold.
 */
const CLEANUP = [
	"create-laju-go", // the installer itself — lives in the repo, not the scaffold
	".devin", // agent skills (graphify, etc.)
	".windsurf", // windsurf rules
	".zero", // agent specialists
	".pi", // pi-lens settings
	"AGENTS.md", // agent instructions (dev-only)
	".mcp.json", // MCP server config (dev-only)
	"site", // Astro docs/landing — repo-only, not part of a scaffold
];

/**
 * Subpaths within .llm-wiki/ to strip — keeps concepts/entities (useful for
 * understanding the repo) but removes session logs, raw captures, and auto-generated metadata.
 */
const WIKI_CLEANUP = [
	".llm-wiki/wiki/sources", // session observations + dev process logs
	".llm-wiki/meta", // auto-generated metadata (regenerable via wiki_rebuild_meta)
	".llm-wiki/raw", // raw source packets (original captures)
	".llm-wiki/.discoveries", // gap analysis (dev artifact)
	".llm-wiki/.obsidian", // Obsidian editor config (dev-only)
];

function detectAvailablePackageManagers() {
	const packageManagers = ["bun", "yarn", "npm"];
	const available = [];

	for (const pm of packageManagers) {
		try {
			execSync(`${pm} --version`, { stdio: "ignore" });
			available.push(pm);
		} catch (e) {}
	}
	return available.length > 0 ? available : ["npm"];
}

async function selectPackageManager(available) {
	if (available.length === 1) {
		return available[0];
	}

	console.log("");
	console.log("\x1b[36m  Which package manager would you like to use?\x1b[0m");

	const response = await prompts({
		type: "select",
		name: "packageManager",
		message: "",
		choices: available.map((pm) => ({
			title:
				pm === "npm"
					? "npm (stable ✓)"
					: pm === "yarn"
						? "yarn (fast ✓)"
						: "bun (experimental ⚠️)",
			value: pm,
		})),
	});

	const selected = response.packageManager || "npm";

	if (selected === "bun") {
		console.log(
			"\x1b[90m  Note: Bun is experimental. If you encounter issues, try npm or yarn.\x1b[0m",
		);
		console.log("");
	}

	return selected;
}

function getInstallCommand(packageManager) {
	switch (packageManager) {
		case "bun":
			return "bun install";
		case "yarn":
			return "yarn";
		case "npm":
		default:
			return "npm install";
	}
}

function getRunCommand(packageManager) {
	switch (packageManager) {
		case "bun":
			return "bun run";
		case "yarn":
			return "yarn";
		case "npm":
		default:
			return "npm run";
	}
}

const ASCII_ART = `
\x1b[38;5;208m
              ++++++++++++++++++
             +++++++++++++++++++
             +++++++++++++++++++
            +++++++++++++++++++
            +++++++++++++++++++
           ++++++++++++++++++++
           +++++++++++++++++++
           +++++++++++++++++++
          +++++++++++++++++++
          +++++++++++++++++++
         +++++++++++++++++++

         +++++++++++++++++++++++++++++++++
        ++++++++++++++++++++++++++++++++++
        +++++++++++++++++++++++++++++++++
       ++++++++++++++++++++++++++++++++++
      ++++++++++++++++++++++++++++++++++
      ++++++++++++++++++++++++++++++++++
      +++++++++++++++++++++++++++++++++
     +++++++++++++++++++++++++++++++++
     ++++++++++++++++++++++++++++++++
\x1b[0m
`;

program
	.name("create-laju-go")
	.description("CLI to create a new Laju Go project from template")
	.version("1.2.0");

program
	.argument("[project-directory]", "Project directory name")
	.option("--package-manager <pm>", "Package manager to use (npm, yarn, bun)")
	.action(async (projectDirectory, options) => {
		try {
			console.log("");
			console.log(ASCII_ART);
			console.log(
				"\x1b[36m  Create a new Go project with Laju Framework\x1b[0m",
			);
			console.log("");

			// Check if Go is installed
			try {
				execSync("go version", { stdio: "ignore" });
			} catch (e) {
				console.log(
					"\x1b[1;31m✖\x1b[0m \x1b[1;91mError:\x1b[0m Go is not installed or not in PATH.",
				);
				console.log("");
				console.log("\x1b[90m  Please install Go to continue:\x1b[0m");
				console.log("  - Download: https://go.dev/dl/");
				console.log("  - macOS: brew install go");
				console.log("  - Linux: sudo apt install golang (Ubuntu/Debian)");
				console.log("");
				process.exit(1);
			}

			let packageManager;
			if (options.packageManager) {
				packageManager = options.packageManager;
				if (!["npm", "yarn", "bun"].includes(packageManager)) {
					console.log(
						"\x1b[1;31m✖\x1b[0m \x1b[1;91mError:\x1b[0m Invalid package manager. Use npm, yarn, or bun.",
					);
					process.exit(1);
				}
			} else {
				const availablePackageManagers = detectAvailablePackageManagers();
				packageManager = await selectPackageManager(availablePackageManagers);
			}
			console.log("\x1b[36m  Using " + packageManager + "\x1b[0m");
			console.log("");

			// If no project name, ask user
			if (!projectDirectory) {
				console.log("\x1b[36m  Project name:\x1b[0m");
				const response = await prompts({
					type: "text",
					name: "projectDirectory",
					message: "",
				});
				projectDirectory = response.projectDirectory;
			}

			if (!projectDirectory) {
				console.log(
					"\x1b[1;31m✖\x1b[0m \x1b[1;91mError:\x1b[0m Project name is required to continue.",
				);
				process.exit(1);
			}

			// Validate project name (npm package name rules)
			const nameRegex =
				/^(?:@[a-z0-9-~][a-z0-9-._~]*\/)?[a-z0-9-~][a-z0-9-._~]*$/;
			if (!nameRegex.test(projectDirectory)) {
				console.log(
					"\x1b[1;31m✖\x1b[0m \x1b[1;91mError:\x1b[0m Invalid project name. Use lowercase letters, numbers, and hyphens only.",
				);
				process.exit(1);
			}

			const targetPath = path.resolve(projectDirectory);

			// Check if directory exists
			if (fs.existsSync(targetPath)) {
				console.log(
					"\x1b[1;31m✖\x1b[0m \x1b[1;91mError:\x1b[0m Directory \x1b[36m" +
						projectDirectory +
						"\x1b[0m already exists. Choose another name.",
				);
				process.exit(1);
			}

			console.log("");
			console.log(
				"\x1b[90m  Creating project at \x1b[36m" + targetPath + "\x1b[0m",
			);
			console.log("");

			// Check if Git is installed
			try {
				execSync("git --version", { stdio: "ignore" });
			} catch (e) {
				console.log(
					"\x1b[1;31m✖\x1b[0m \x1b[1;91mError:\x1b[0m Git is not installed or not in PATH.",
				);
				console.log("");
				console.log("\x1b[90m  Please install Git to continue:\x1b[0m");
				console.log("  - Windows: https://git-scm.com/download/win");
				console.log("  - macOS: brew install git");
				console.log("  - Linux: sudo apt install git (Ubuntu/Debian)");
				console.log("");
				process.exit(1);
			}

			// Clone template from GitHub
			const emitter = degit("maulanashalihin/laju-go");

			await emitter.clone(targetPath);

			// Strip dev-only / repo-internal files from the scaffolded project.
			console.log("\x1b[36m  Cleaning up dev-only files...\x1b[0m");
			for (const rel of [...CLEANUP, ...WIKI_CLEANUP]) {
				const fullPath = path.join(targetPath, rel);
				fs.rmSync(fullPath, { recursive: true, force: true });
			}
			console.log("\x1b[32m  ✓ Dev-only files removed\x1b[0m");
			console.log("");

			// Read package.json from template
			const packageJsonPath = path.join(targetPath, "package.json");
			const packageJson = require(packageJsonPath);

			// Update project name in package.json
			packageJson.name = projectDirectory;

			// Reset version to 0.0.1 for new project
			packageJson.version = "0.0.1";

			// Write back package.json
			fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2));

			// Change directory and run setup commands
			const originalDir = process.cwd();
			process.chdir(targetPath);

			try {
				console.log("\x1b[36m  Installing Go dependencies...\x1b[0m");
				console.log("");
				execSync("go mod download", { stdio: "inherit", timeout: 120000 });
				console.log("\x1b[32m  ✓ Go dependencies installed\x1b[0m");
				console.log("");

				console.log("\x1b[36m  Installing Node.js dependencies...\x1b[0m");
				console.log("");
				execSync(getInstallCommand(packageManager), {
					stdio: "inherit",
					timeout: 300000,
				});
				console.log("\x1b[32m  ✓ Node.js dependencies installed\x1b[0m");
				console.log("");

				console.log("\x1b[36m  Setting up environment...\x1b[0m");
				execSync(
					process.platform === "win32"
						? "copy .env.example .env"
						: "cp .env.example .env",
					{ stdio: "pipe", timeout: 10000 },
				);
				console.log("\x1b[32m  ✓ Environment configured\x1b[0m");
				console.log("");
			} finally {
				process.chdir(originalDir);
			}

			console.log("");
			console.log("\x1b[1;36m  ✓ Project created successfully!\x1b[0m");
			console.log("");
			console.log("\x1b[90m  Next steps:\x1b[0m");
			console.log("");
			console.log("    cd " + projectDirectory);
			console.log("    " + getRunCommand(packageManager) + " dev:all");
			console.log("");
			console.log("\x1b[90m  Learn more: https://laju.dev\x1b[0m");
			console.log(
				"\x1b[90m  Docs: https://github.com/maulanashalihin/laju-go\x1b[0m",
			);
			console.log("");
		} catch (error) {
			console.error("Error:", error.message);
			if (process.env.DEBUG) {
				console.error(error.stack);
			}
			process.exit(1);
		}
	});

program.parse();
