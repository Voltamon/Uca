import { vi } from 'vitest';

// Mock fetch globally to avoid "Invalid URL" errors in Node/JSDOM
// and to allow components to "mount" even if they call APIs on init.
vi.stubGlobal('fetch', vi.fn(() => 
  Promise.resolve({
    ok: true,
    status: 200,
    json: () => Promise.resolve([]),
  } as Response)
));
