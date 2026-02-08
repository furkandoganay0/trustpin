const strings = {
  tr: {
    "ui.back": "TrustPin’a Dön",
    "ui.search": "Dokümanlarda ara",
    "sidebar.gettingStarted": "Başlarken",
    "sidebar.intro": "Giriş",
    "sidebar.quickstart": "Hızlı Başlangıç",
    "sidebar.api": "API",
    "sidebar.apiOverview": "API Genel Bakış",
    "sidebar.openapi": "OpenAPI Referansı",
    "sidebar.security": "Güvenlik",
    "sidebar.securityModel": "Güvenlik Modeli",
    "sidebar.deploy": "Kurulum",
    "sidebar.support": "Destek",
    "sidebar.faq": "SSS",
    "intro.badge": "Giriş",
    "intro.title": "TrustPin – Banka seviyesinde mobil onay & MFA backend",
    "intro.p1": "TrustPin güvenlik odaklı, kurum içi kurulabilen bir onay sistemidir. Kullanıcı doğrulamaz; sadece kriptografik onayları doğrular. Bankalar, fintech’ler ve kurumlar için açık kaynak ve self‑hosted olarak tasarlanmıştır.",
    "intro.calloutTitle": "TrustPin nedir",
    "intro.calloutBody": "Katı durum makineleri, append‑only denetim logları ve cihaz imzalı challenge’larla çalışan kriptografik bir onay altyapısıdır.",
    "intro.notTitle": "TrustPin ne değildir",
    "intro.notBody": "Kimlik sağlayıcınızın yerini almaz. Kimlik doğrulaması mevcut sisteminizde kalır.",
    "intro.conceptsTitle": "Temel kavramlar",
    "intro.conceptsBody": "Deterministik payload ve imzalı challenge’lar doğrulanabilir onay sağlar. Sunucu yalnızca public key saklar.",
    "intro.archTitle": "Üst seviye mimari",
    "intro.archBody": "API stateless çalışır; PostgreSQL gerçek durum kaynağıdır, Redis TTL cache ve nonce tekilliği için kullanılır, push sağlayıcı cihazlara bildirim gönderir.",
    "intro.list1": "API servisi (Go) ve katı HTTP handler’lar",
    "intro.list2": "Durum makineleri ve audit log için PostgreSQL",
    "intro.list3": "TTL cache, rate limit ve nonce tekilliği için Redis",
    "intro.list4": "Push provider soyutlaması (mock edilebilir)",
    "intro.flowTitle": "Tipik akış",
    "intro.flowBody": "Enrollment → Device activation → Challenge initiation → Approve/Reject → Audit event.",
    "quick.badge": "Hızlı Başlangıç",
    "quick.title": "TrustPin’i yerelde çalıştırın",
    "quick.p1": "Docker Compose ile API, PostgreSQL ve Redis’i ayağa kaldırın.",
    "quick.step1": "1. Stack’i başlatın",
    "quick.step2": "2. API canlılığını doğrulayın",
    "quick.step2Body": "Varsayılan olarak API şu adreste çalışır: http://localhost:8081",
    "quick.step3": "3. Dokümanları açın",
    "quick.docs1": "Geliştirici dokümanları:",
    "quick.docs2": "Swagger/OpenAPI:",
    "quick.nextTitle": "Sonraki adımlar",
    "quick.nextBody": "Enrollment ile başlayın, cihazı aktif edin, ardından challenge akışını tamamlayın.",
    "api.badge": "API Genel Bakış",
    "api.title": "TrustPin API",
    "api.p1": "TrustPin enrollment, cihaz yönetimi, challenge yaşam döngüsü ve audit erişimi için minimal endpoint seti sunar.",
    "api.baseTitle": "Base URL",
    "api.coreTitle": "Temel endpoint’ler",
    "api.e1": "/v1/enrollments/init — enrollment ve pairing code oluşturur",
    "api.e2": "/v1/devices/activate — public key ile cihazı aktive eder",
    "api.e3": "/v1/auth/challenges/init — challenge başlatır",
    "api.e4": "/v1/auth/challenges/{id} — canonical payload getirir",
    "api.e5": "/v1/auth/challenges/{id}/approve — imza ile onaylar",
    "api.e6": "/v1/auth/challenges/{id}/reject — imza ile reddeder",
    "api.e7": "/v1/auth/challenges/{id}/status — durum sorgular",
    "api.e8": "/v1/audit/events — audit kayıtlarını listeler",
    "api.ex1": "Örnek: enrollment başlat",
    "api.ex2": "Örnek: challenge onayla",
    "api.ex2Body": "Onay için cihaz imzası ve canonical payload gereklidir.",
    "api.openapiTitle": "OpenAPI referansı",
    "api.openapiBody": "İstek/yanıt şemaları için Swagger UI’yi kullanın.",
    "api.openapiCta": "OpenAPI Dokümanları",
    "sec.badge": "Güvenlik",
    "sec.title": "Güvenlik modeli",
    "sec.p1": "TrustPin sıfır güven ortamları için tasarlanmıştır. Her onay kriptografik olarak doğrulanır.",
    "sec.pkTitle": "Sadece public key saklama",
    "sec.pkBody": "Cihazlar Ed25519 anahtar çiftini yerelde üretir. Sunucu yalnızca public key saklar.",
    "sec.canonTitle": "Canonical payload’lar",
    "sec.canonBody": "Deterministik payload sıralaması, imza karışıklığını ve replay saldırılarını engeller.",
    "sec.nonceTitle": "Nonce tekilliği + TTL",
    "sec.nonceBody": "Her challenge benzersiz nonce içerir ve Redis süresince tekrar kullanımını engeller.",
    "sec.calloutTitle": "Garanti",
    "sec.calloutBody": "Veritabanı erişimi olsa bile, cihaz imzası olmadan onay verilemez.",
    "sec.auditTitle": "Audit log",
    "sec.auditBody": "Tüm geçişler append‑only audit event üretir.",
    "dep.badge": "Kurulum",
    "dep.title": "TrustPin’i kurum içinde dağıtın",
    "dep.p1": "TrustPin kendi altyapınızda çalışacak şekilde tasarlanmıştır. Varsayılan compose prod için bir başlangıçtır.",
    "dep.coreTitle": "Temel servisler",
    "dep.core1": "API servisi (stateless)",
    "dep.core2": "PostgreSQL (gerçek durum kaynağı)",
    "dep.core3": "Redis (TTL cache, nonce tekilliği)",
    "dep.checkTitle": "Prod kontrol listesi",
    "dep.check1": "TLS sonlandırma (ingress veya reverse proxy) etkinleştir",
    "dep.check2": "Güçlü veritabanı parolaları kullan ve secret’ları döndür",
    "dep.check3": "PostgreSQL ve audit verilerini düzenli yedekle",
    "dep.check4": "Docker image tag’lerini sabitle",
    "dep.check5": "Ağ erişimini iç zonlarla sınırla",
    "dep.cfgTitle": "Konfigürasyon",
    "dep.cfgBody": "Ortam değişkenleriyle yapılandırılır. TTL ve rate limit değerlerini compose üzerinden ayarlayın.",
    "faq.badge": "SSS",
    "faq.title": "Sık sorulan sorular",
    "faq.q1": "TrustPin bir kimlik doğrulama sağlayıcısı mı?",
    "faq.a1": "Hayır. TrustPin yalnızca kriptografik onayları doğrular.",
    "faq.q2": "Self‑hosted çalıştırabilir miyim?",
    "faq.a2": "Evet. TrustPin açık kaynak ve kurum içi kurulum için tasarlanmıştır.",
    "faq.q3": "Private key saklıyor mu?",
    "faq.a3": "Hayır. Private key cihazdan çıkmaz; sunucu yalnızca public key saklar.",
    "faq.q4": "API şemalarını nerede görebilirim?",
    "faq.a4": "Etkileşimli OpenAPI referansı /swagger/ altında.",
    "toc.title": "Bu sayfada"
  },
  en: {
    "ui.back": "Back to TrustPin",
    "ui.search": "Search docs",
    "sidebar.gettingStarted": "Getting Started",
    "sidebar.intro": "Introduction",
    "sidebar.quickstart": "Quickstart",
    "sidebar.api": "API",
    "sidebar.apiOverview": "API Overview",
    "sidebar.openapi": "OpenAPI Reference",
    "sidebar.security": "Security",
    "sidebar.securityModel": "Security Model",
    "sidebar.deploy": "Deployment",
    "sidebar.support": "Support",
    "sidebar.faq": "FAQ",
    "intro.badge": "Introduction",
    "intro.title": "TrustPin – Bank‑grade mobile approval & MFA backend",
    "intro.p1": "TrustPin is a security‑first, on‑prem deployable approval system. It does not authenticate users; it only verifies cryptographic approval. Built for banks, fintechs, and enterprises, TrustPin is open‑source and self‑hostable.",
    "intro.calloutTitle": "What TrustPin is",
    "intro.calloutBody": "A cryptographic approval backend with strict state machines, append‑only audit trails, and device‑signed challenges.",
    "intro.notTitle": "What TrustPin is not",
    "intro.notBody": "TrustPin does not replace your identity provider. Authentication stays in your existing stack.",
    "intro.conceptsTitle": "Key concepts",
    "intro.conceptsBody": "Deterministic payloads and signed challenges guarantee verifiable approvals. The server stores only public keys.",
    "intro.archTitle": "High‑level architecture",
    "intro.archBody": "The API runs statelessly; PostgreSQL is the source of truth, Redis enforces TTL caches and nonce uniqueness, and push providers notify devices.",
    "intro.list1": "API service (Go) with strict HTTP handlers",
    "intro.list2": "PostgreSQL for state machines and audit log",
    "intro.list3": "Redis for TTL caches, rate limits, and nonce uniqueness",
    "intro.list4": "Push provider abstraction (mockable)",
    "intro.flowTitle": "Typical flow",
    "intro.flowBody": "Enrollment → Device activation → Challenge initiation → Approve/Reject → Audit event.",
    "quick.badge": "Quickstart",
    "quick.title": "Run TrustPin locally",
    "quick.p1": "Use Docker Compose to start the API, PostgreSQL, and Redis.",
    "quick.step1": "1. Start the stack",
    "quick.step2": "2. Verify the API is live",
    "quick.step2Body": "The API listens on http://localhost:8081 by default.",
    "quick.step3": "3. Open the docs",
    "quick.docs1": "Developer docs:",
    "quick.docs2": "Swagger/OpenAPI:",
    "quick.nextTitle": "Next steps",
    "quick.nextBody": "Start with enrollment, activate a device, then create a challenge end‑to‑end.",
    "api.badge": "API Overview",
    "api.title": "TrustPin API",
    "api.p1": "TrustPin exposes a minimal set of endpoints for enrollment, device management, challenge lifecycle, and audit access.",
    "api.baseTitle": "Base URL",
    "api.coreTitle": "Core endpoints",
    "api.e1": "/v1/enrollments/init — create enrollment and pairing code",
    "api.e2": "/v1/devices/activate — activate device with public key",
    "api.e3": "/v1/auth/challenges/init — issue challenge",
    "api.e4": "/v1/auth/challenges/{id} — fetch canonical payload",
    "api.e5": "/v1/auth/challenges/{id}/approve — approve with signature",
    "api.e6": "/v1/auth/challenges/{id}/reject — reject with signature",
    "api.e7": "/v1/auth/challenges/{id}/status — query state",
    "api.e8": "/v1/audit/events — list audit events",
    "api.ex1": "Example: initiate enrollment",
    "api.ex2": "Example: approve challenge",
    "api.ex2Body": "Approval requires the device signature and canonical payload.",
    "api.openapiTitle": "OpenAPI reference",
    "api.openapiBody": "Use the interactive Swagger UI for request/response schemas.",
    "api.openapiCta": "OpenAPI Docs",
    "sec.badge": "Security",
    "sec.title": "Security model",
    "sec.p1": "TrustPin is designed for zero‑trust environments. Every approval must be cryptographically verifiable.",
    "sec.pkTitle": "Public‑key only storage",
    "sec.pkBody": "Devices generate Ed25519 keypairs locally. The server stores only public keys.",
    "sec.canonTitle": "Canonical payloads",
    "sec.canonBody": "Deterministic payload ordering prevents signature confusion and replay attacks.",
    "sec.nonceTitle": "Nonce uniqueness + TTLs",
    "sec.nonceBody": "Each challenge includes a unique nonce. Redis enforces uniqueness until expiration.",
    "sec.calloutTitle": "Guarantee",
    "sec.calloutBody": "Even with database access, an attacker cannot approve without the device signature.",
    "sec.auditTitle": "Audit log",
    "sec.auditBody": "All transitions produce append‑only audit events.",
    "dep.badge": "Deployment",
    "dep.title": "Deploy TrustPin on‑prem",
    "dep.p1": "TrustPin is designed to run inside your own infrastructure. The default docker‑compose stack is a starting point for production hardening.",
    "dep.coreTitle": "Core services",
    "dep.core1": "API service (stateless)",
    "dep.core2": "PostgreSQL (source of truth)",
    "dep.core3": "Redis (TTL cache, nonce uniqueness)",
    "dep.checkTitle": "Production checklist",
    "dep.check1": "Enable TLS termination (ingress or reverse proxy)",
    "dep.check2": "Set strong database credentials and rotate secrets",
    "dep.check3": "Back up PostgreSQL and audit data regularly",
    "dep.check4": "Pin docker image tags",
    "dep.check5": "Restrict network access to internal zones",
    "dep.cfgTitle": "Configuration",
    "dep.cfgBody": "Configure via environment variables. Adjust TTLs and rate limits in compose.",
    "faq.badge": "FAQ",
    "faq.title": "Frequently asked questions",
    "faq.q1": "Is TrustPin an authentication provider?",
    "faq.a1": "No. TrustPin only verifies approvals via cryptographic signatures.",
    "faq.q2": "Can I self‑host?",
    "faq.a2": "Yes. TrustPin is open‑source and built for on‑prem deployment.",
    "faq.q3": "Does TrustPin store private keys?",
    "faq.a3": "No. Private keys never leave the device. TrustPin stores only public keys.",
    "faq.q4": "Where do I see API schemas?",
    "faq.a4": "The interactive OpenAPI reference is available at /swagger/.",
    "toc.title": "On this page"
  }
};

const applyLanguage = (lang) => {
  document.documentElement.lang = lang;
  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const key = el.getAttribute("data-i18n");
    if (strings[lang] && strings[lang][key]) {
      el.textContent = strings[lang][key];
    }
  });
  const input = document.querySelector("#doc-search");
  if (input) {
    input.setAttribute("placeholder", strings[lang]["ui.search"]);
  }
  const btn = document.querySelector("#lang-toggle");
  if (btn) {
    btn.textContent = lang === "tr" ? "EN" : "TR";
    btn.setAttribute("aria-label", lang === "tr" ? "Switch to English" : "Türkçeye geç");
  }
  localStorage.setItem("trustpin_docs_lang", lang);
};

const highlightActiveLink = () => {
  const path = window.location.pathname.replace(/\/$/, "");
  document.querySelectorAll(".sidebar a").forEach((link) => {
    const target = link.getAttribute("href").replace(/\/$/, "");
    if (target && path.endsWith(target)) {
      link.classList.add("active");
    }
  });
};

const buildToc = () => {
  const toc = document.querySelector(".toc");
  const content = document.querySelector(".content");
  if (!toc || !content) return;
  const headings = content.querySelectorAll("h2, h3");
  if (!headings.length) return;
  const list = document.createElement("div");
  headings.forEach((heading) => {
    if (!heading.id) {
      heading.id = heading.textContent.toLowerCase().replace(/[^a-z0-9]+/g, "-");
    }
    const link = document.createElement("a");
    link.href = `#${heading.id}`;
    link.textContent = heading.textContent;
    list.appendChild(link);
  });
  toc.appendChild(list);
};

const bindSearch = () => {
  const input = document.querySelector("#doc-search");
  const links = Array.from(document.querySelectorAll(".sidebar a"));
  if (!input) return;
  input.addEventListener("input", () => {
    const q = input.value.toLowerCase();
    links.forEach((link) => {
      const text = link.textContent.toLowerCase();
      link.style.display = text.includes(q) ? "block" : "none";
    });
  });
};

const savedLang = localStorage.getItem("trustpin_docs_lang") || "tr";
applyLanguage(savedLang);
const toggle = document.querySelector("#lang-toggle");
if (toggle) {
  toggle.addEventListener("click", () => {
    const next = document.documentElement.lang === "tr" ? "en" : "tr";
    applyLanguage(next);
  });
}

highlightActiveLink();
buildToc();
bindSearch();
