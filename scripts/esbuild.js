const childProcess = require('child_process')
const path = require('path')
const fs = require('fs')
const os = require('os')

const repoDir = path.dirname(__dirname)
const npmDir = path.join(repoDir, 'npm', 'cg')
const version = fs.readFileSync(path.join(repoDir, 'version.txt'), 'utf8').trim()
const nodeTarget = 'node10'; // See: https://nodejs.org/en/about/releases/

const buildNeutralLib = (esbuildPath) => {
  const binDir = path.join(npmDir, 'bin')
  fs.mkdirSync(binDir, { recursive: true })

  // Generate "npm/cg/install.js"
  childProcess.execFileSync(esbuildPath, [
    path.join(repoDir, 'lib', 'npm', 'node-install.ts'),
    '--outfile=' + path.join(npmDir, 'install.js'),
    '--bundle',
    '--target=' + nodeTarget,
    // Note: https://socket.dev have complained that inlining the version into
    // the install script messes up some internal scanning that they do by
    // making it seem like cg's install script code changes with every
    // cg release. So now we read it from "package.json" instead.
    // '--define:CG_VERSION=' + JSON.stringify(version),
    '--external:cg',
    '--platform=node',
    '--log-level=warning',
  ], { cwd: repoDir })

  // Generate "npm/cg/bin/cg"
  childProcess.execFileSync(esbuildPath, [
    path.join(repoDir, 'lib', 'npm', 'node-shim.ts'),
    '--outfile=' + path.join(binDir, 'cg'),
    '--bundle',
    '--target=' + nodeTarget,
    '--define:CG_VERSION=' + JSON.stringify(version),
    '--external:cg',
    '--platform=node',
    '--log-level=warning',
  ], { cwd: repoDir })

  // Get supported platforms
  const platforms = { exports: {} }
  new Function('module', 'exports', 'require', childProcess.execFileSync(esbuildPath, [
    path.join(repoDir, 'lib', 'npm', 'node-platform.ts'),
    '--bundle',
    '--target=' + nodeTarget,
    '--external:cg',
    '--platform=node',
    '--log-level=warning',
  ], { cwd: repoDir }))(platforms, platforms.exports, require)
  const optionalDependencies = Object.fromEntries(Object.values({
    ...platforms.exports.knownWindowsPackages,
    ...platforms.exports.knownUnixlikePackages,
    ...platforms.exports.knownWebAssemblyFallbackPackages,
  }).sort().map(x => [x, version]))

  // Update "npm/cg/package.json"
  const pjPath = path.join(npmDir, 'package.json')
  const package_json = JSON.parse(fs.readFileSync(pjPath, 'utf8'))
  package_json.optionalDependencies = optionalDependencies
  fs.writeFileSync(pjPath, JSON.stringify(package_json, null, 2) + '\n')
}

const updateVersionPackageJSON = pathToPackageJSON => {
  const version = fs.readFileSync(path.join(path.dirname(__dirname), 'version.txt'), 'utf8').trim()
  const json = JSON.parse(fs.readFileSync(pathToPackageJSON, 'utf8'))

  if (json.version !== version) {
    json.version = version
    fs.writeFileSync(pathToPackageJSON, JSON.stringify(json, null, 2) + '\n')
  }
}

const updateVersionGo = () => {
  const version_txt = fs.readFileSync(path.join(repoDir, 'version.txt'), 'utf8').trim()
  const version_go = `package main\n\nconst cgVersion = "${version_txt}"\n`
  const version_go_path = path.join(repoDir, 'cmd', 'cg', 'version.go')

  // Update this atomically to avoid issues with this being overwritten during use
  const temp_path = version_go_path + Math.random().toString(36).slice(1)
  fs.writeFileSync(temp_path, version_go)
  fs.renameSync(temp_path, version_go_path)
}

// This is helpful for ES6 modules which don't have access to __dirname
exports.dirname = __dirname

// The main Makefile invokes this script before publishing
if (require.main === module) {
  if (process.argv.indexOf('--version') >= 0) updateVersionPackageJSON(process.argv[2])
  else if (process.argv.indexOf('--neutral') >= 0) buildNeutralLib(process.argv[2])
  else if (process.argv.indexOf('--update-version-go') >= 0) updateVersionGo()
  else throw new Error('Expected a flag')
}
