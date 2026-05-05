import { createServer } from 'node:http'
import { readFile } from 'node:fs/promises'
import { join, extname } from 'node:path'
import { fileURLToPath } from 'node:url'
import http from 'node:http'

const STATIC_DIR = join(fileURLToPath(new URL('./.output/public', import.meta.url)))
const BACKEND = 'http://localhost:8080'
const PORT = 3000

const MIME = {
  '.html': 'text/html; charset=utf-8',
  '.css': 'text/css',
  '.js': 'application/javascript',
  '.mjs': 'application/javascript',
  '.json': 'application/json',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.svg': 'image/svg+xml',
  '.ico': 'image/x-icon',
}

function proxyReq(req, res) {
  const opts = new URL(req.url, BACKEND)
  const proxy = http.request(opts.href, { method: req.method, headers: req.headers }, (pres) => {
    res.writeHead(pres.statusCode, pres.headers)
    pres.pipe(res)
  })
  proxy.on('error', () => { res.writeHead(502); res.end('Bad Gateway') })
  req.pipe(proxy)
}

function serveStatic(req, res) {
  let path = req.url.split('?')[0]
  if (path.endsWith('/')) path += 'index.html'
  const file = join(STATIC_DIR, path)
  readFile(file).then((data) => {
    res.writeHead(200, { 'Content-Type': MIME[extname(file)] || 'application/octet-stream' })
    res.end(data)
  }).catch(() => {
    // SPA fallback: serve index.html for non-file routes
    readFile(join(STATIC_DIR, 'index.html')).then((data) => {
      res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' })
      res.end(data)
    }).catch(() => { res.writeHead(404); res.end('Not Found') })
  })
}

createServer((req, res) => {
  if (req.url.startsWith('/api/')) return proxyReq(req, res)
  serveStatic(req, res)
}).listen(PORT, () => {
  console.log(`Frontend: http://localhost:${PORT}`)
  console.log(`API proxy: ${BACKEND}`)
})
