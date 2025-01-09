const path = require('path');
const fs = require('fs');

const FOLDERS = [
  'frontend/apps',
  'frontend/libs'
]

const ROOT_DIR = path.resolve(__dirname, '..');

/**
 * @type {{ promise: Promise<{ default: unknown }>, rootPath: string }[]}
 */
const configPromises = []

FOLDERS.forEach(folder => {
  const fullPath = path.resolve(ROOT_DIR, folder);
  fs.readdirSync(fullPath).forEach(file => {
    const rootPath = path.resolve(fullPath, file);
    const eslintPath = path.resolve(rootPath, 'eslint.config.js');
    if (fs.existsSync(eslintPath)) {
      configPromises.push({
        promise: import(eslintPath),
        rootPath
      });
    }
  })
})

let idx = 0;

function next() {
  if (idx >= configPromises.length) {
    return;
  }

  configPromises[idx].promise.then(config => {
    runEslint(config.default, configPromises[idx].rootPath).then(() => {
      idx++
      next()
    });
  })
}

async function runEslint(config, rootPath) {
  const eslint = require('eslint');

  const cli = new eslint.ESLint({
    overrideConfig: config,
    overrideConfigFile: true,
    cwd: ROOT_DIR,
    fix: true,
    // errorOnUnmatchedPattern: false,
  });

  console.log(`Running eslint on ${rootPath}`);

  const report = await cli.lintFiles([
    path.resolve(rootPath, 'src/lib/*.tsx'),
  ]);


  if (report.errorCount > 0) {
    console.error(cli.getFormatter()(report.results));
    process.exit(1);
  } else {
    console.debug(`No errors found in ${rootPath}, applying changes...`);
  }

  const now = performance.now();

  for (const file of report) {
    if (!file.output) continue;

    fs.writeFileSync(file.filePath, file.output)
  }

  console.log(`Done in ${performance.now() - now}ms`);
}

next()
