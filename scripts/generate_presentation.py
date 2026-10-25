"""
CustardFin Presentation Generator
Generates a modern, widescreen (16:9), dark-themed PowerPoint presentation (.pptx)
with native OpenXML slide transitions, shape build animations, cards, and diagrams.
Also exports the presentation XML for reference and inspection.
"""

import os
import sys
import zipfile
import xml.etree.ElementTree as ET
from xml.dom import minidom

import pptx
from pptx import Presentation
from pptx.util import Inches, Pt
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN
from pptx.enum.shapes import MSO_SHAPE
from pptx.oxml import parse_xml
from pptx.oxml.ns import nsdecls

# Palette Tokens
COLOR_BG = RGBColor(9, 13, 22)        # #090D16
COLOR_BG_CARD = RGBColor(15, 23, 42)   # #0F172A
COLOR_CARD_BOX = RGBColor(30, 41, 59)  # #1E293B
COLOR_BORDER = RGBColor(71, 85, 105)   # #475569
COLOR_TEXT_MAIN = RGBColor(248, 250, 252) # #F8FAFC
COLOR_TEXT_MUTED = RGBColor(148, 163, 184) # #94A3B8
COLOR_TEXT_DIM = RGBColor(100, 116, 139)   # #64748B
COLOR_GOLD = RGBColor(245, 158, 11)    # #F59E0B
COLOR_EMERALD = RGBColor(16, 185, 129) # #10B981
COLOR_CYAN = RGBColor(6, 182, 212)     # #06B6D4
COLOR_RED = RGBColor(239, 68, 68)      # #EF4444

def create_presentation():
    prs = Presentation()
    prs.slide_width = Inches(13.333)
    prs.slide_height = Inches(7.5)
    blank_layout = prs.slide_layouts[6]

    def add_slide_base(transition="fade", push_dir="l"):
        slide = prs.slides.add_slide(blank_layout)
        
        # OpenXML Transition
        if transition == "push":
            trans_xml = f'<p:transition {nsdecls("p")} spd="med"><p:push dir="{push_dir}"/></p:transition>'
        elif transition == "wipe":
            trans_xml = f'<p:transition {nsdecls("p")} spd="med"><p:wipe dir="r"/></p:transition>'
        else:
            trans_xml = f'<p:transition {nsdecls("p")} spd="med"><p:fade filter="out"/></p:transition>'
        slide.element.append(parse_xml(trans_xml))

        # Solid Dark Background Shape
        bg = slide.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, 0, Inches(13.333), Inches(7.5))
        bg.fill.solid()
        bg.fill.fore_color.rgb = COLOR_BG
        bg.line.fill.background()
        
        # Subtle Top Ambient Bar
        top_bar = slide.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, 0, Inches(13.333), Inches(0.08))
        top_bar.fill.solid()
        top_bar.fill.fore_color.rgb = COLOR_GOLD
        top_bar.line.fill.background()

        return slide

    def add_header(slide, title, category="CUSTARDFIN • DISTRIBUTED COMMUNITY MEDIA"):
        # Category Eyebrow
        tb_cat = slide.shapes.add_textbox(Inches(0.8), Inches(0.45), Inches(11.7), Inches(0.35))
        tf_cat = tb_cat.text_frame
        tf_cat.word_wrap = True
        p_cat = tf_cat.paragraphs[0]
        p_cat.text = category.upper()
        p_cat.font.size = Pt(10)
        p_cat.font.bold = True
        p_cat.font.color.rgb = COLOR_GOLD

        # Slide Main Title
        tb_title = slide.shapes.add_textbox(Inches(0.8), Inches(0.75), Inches(11.7), Inches(0.7))
        tf_title = tb_title.text_frame
        tf_title.word_wrap = True
        p_title = tf_title.paragraphs[0]
        p_title.text = title
        p_title.font.size = Pt(24)
        p_title.font.bold = True
        p_title.font.color.rgb = COLOR_TEXT_MAIN

        # Subtle Divider
        divider = slide.shapes.add_shape(MSO_SHAPE.RECTANGLE, Inches(0.8), Inches(1.5), Inches(11.733), Inches(0.02))
        divider.fill.solid()
        divider.fill.fore_color.rgb = RGBColor(30, 41, 59)
        divider.line.fill.background()

    def add_card(slide, x, y, w, h, title, items, tag=None, tag_color=COLOR_GOLD, border_color=COLOR_BORDER):
        # Card Background
        card = slide.shapes.add_shape(MSO_SHAPE.ROUNDED_RECTANGLE, Inches(x), Inches(y), Inches(w), Inches(h))
        card.fill.solid()
        card.fill.fore_color.rgb = COLOR_BG_CARD
        card.line.color.rgb = border_color
        card.line.width = Pt(1.2)

        # Content Text Box
        pad_x = 0.25
        pad_y = 0.2
        tb = slide.shapes.add_textbox(Inches(x + pad_x), Inches(y + pad_y), Inches(w - 2 * pad_x), Inches(h - 2 * pad_y))
        tf = tb.text_frame
        tf.word_wrap = True

        p_first = tf.paragraphs[0]
        if tag:
            p_first.text = f"[{tag}] {title}"
        else:
            p_first.text = title
        p_first.font.size = Pt(14)
        p_first.font.bold = True
        p_first.font.color.rgb = tag_color
        p_first.space_after = Pt(8)

        for item in items:
            p = tf.add_paragraph()
            p.text = f"•  {item}"
            p.font.size = Pt(11)
            p.font.color.rgb = COLOR_TEXT_MUTED
            p.space_after = Pt(5)

        return card.shape_id

    def inject_shape_animations(slide, shape_ids):
        """Builds OpenXML p:timing main sequence with staggered fade-in entrance effects."""
        if not shape_ids:
            return

        child_nodes = []
        c_id = 3
        for i, spid in enumerate(shape_ids):
            node_type = "clickEffect" if i == 0 else "afterEffect"
            delay = 0 if i == 0 else 250
            node_xml = f'''
            <p:par>
              <p:cTn id="{c_id}" fill="hold">
                <p:stCondLst><p:cond delay="0"/></p:stCondLst>
                <p:childTnLst>
                  <p:par>
                    <p:cTn id="{c_id+1}" fill="hold">
                      <p:stCondLst><p:cond delay="0"/></p:stCondLst>
                      <p:childTnLst>
                        <p:par>
                          <p:cTn id="{c_id+2}" presetID="1" presetClass="entr" presetSubtype="0" fill="hold" nodeType="{node_type}">
                            <p:stCondLst><p:cond delay="{delay}"/></p:stCondLst>
                            <p:childTnLst>
                              <p:set>
                                <p:cBhvr>
                                  <p:cTn id="{c_id+3}" dur="1" fill="hold"/>
                                  <p:tgtEl><p:spTgt spid="{spid}"/></p:tgtEl>
                                  <p:attrNameLst><p:attrName>style.visibility</p:attrName></p:attrNameLst>
                                </p:cBhvr>
                                <p:to><p:strVal val="visible"/></p:to>
                              </p:set>
                              <p:animEffect transition="in" filter="fade">
                                <p:cBhvr>
                                  <p:cTn id="{c_id+4}" dur="400"/>
                                  <p:tgtEl><p:spTgt spid="{spid}"/></p:tgtEl>
                                </p:cBhvr>
                              </p:animEffect>
                            </p:childTnLst>
                          </p:cTn>
                        </p:par>
                      </p:childTnLst>
                    </p:cTn>
                  </p:par>
                </p:childTnLst>
              </p:cTn>
            </p:par>
            '''
            child_nodes.append(node_xml)
            c_id += 5

        timing_xml = f'''
        <p:timing {nsdecls("p")}>
          <p:tnLst>
            <p:par>
              <p:cTn id="1" dur="indefinite" restart="always" nodeType="tmRoot">
                <p:childTnLst>
                  <p:seq concurrent="1" nextAc="seek">
                    <p:cTn id="2" dur="indefinite" nodeType="mainSeq">
                      <p:childTnLst>
                        {''.join(child_nodes)}
                      </p:childTnLst>
                    </p:cTn>
                  </p:seq>
                </p:childTnLst>
              </p:cTn>
            </p:par>
          </p:tnLst>
        </p:timing>
        '''
        slide.element.append(parse_xml(timing_xml))

    # =========================================================================
    # SLIDE 1: TITLE SLIDE (Hero)
    # =========================================================================
    s1 = add_slide_base(transition="fade")
    
    # Hero Title Box
    tb1 = s1.shapes.add_textbox(Inches(1.2), Inches(1.8), Inches(11.0), Inches(2.2))
    tf1 = tb1.text_frame
    tf1.word_wrap = True
    
    p = tf1.paragraphs[0]
    p.text = "CUSTARDFIN"
    p.font.size = Pt(46)
    p.font.bold = True
    p.font.color.rgb = COLOR_GOLD

    p_sub = tf1.add_paragraph()
    p_sub.text = "Distributed Community Media Platform • Hybrid Edge-Cloud Architecture"
    p_sub.font.size = Pt(20)
    p_sub.font.bold = True
    p_sub.font.color.rgb = COLOR_TEXT_MAIN
    p_sub.space_before = Pt(8)

    p_desc = tf1.add_paragraph()
    p_desc.text = "Personal media streaming (Jellyfin/Plex style) meets community discovery (YouTube style) with Zero Cloud Video Storage."
    p_desc.font.size = Pt(13)
    p_desc.font.color.rgb = COLOR_TEXT_MUTED
    p_desc.space_before = Pt(12)

    # Meta Info Cards
    card1 = add_card(s1, 1.2, 4.5, 3.4, 2.2, "Team Credits", [
        "Anish Ghadi (25P0320009)",
        "Anosh Chodankar (25P0320010)",
        "MCA Cloud Computing Lab Project",
        "Semester End Evaluation 2026"
    ], tag="AUTHORS", tag_color=COLOR_CYAN)

    card2 = add_card(s1, 5.0, 4.5, 3.4, 2.2, "Core Philosophy", [
        "Invariant 1: Zero cloud video binaries",
        "Home stays the sovereign origin",
        "Public sharing = ~3 KB JSON manifest",
        "100% Free & Open Source stack"
    ], tag="INVARIANTS", tag_color=COLOR_GOLD)

    card3 = add_card(s1, 8.8, 4.5, 3.4, 2.2, "Demonstration Stack", [
        "Go Backend + Next.js 15 Web UI",
        "SQLite Edge + PostgreSQL Cloud",
        "Hardware NVENC 4s HLS Transcoding",
        "Origin Shield Pull-Through Cache"
    ], tag="LIVE DEMO", tag_color=COLOR_EMERALD)

    inject_shape_animations(s1, [card1, card2, card3])

    # =========================================================================
    # SLIDE 2: THE PROBLEM & MOTIVATION
    # =========================================================================
    s2 = add_slide_base(transition="push", push_dir="l")
    add_header(s2, "The Core Problem: Centralized Clouds vs Isolated Silos", "MOTIVATION & BACKGROUND")

    c1 = add_card(s2, 0.8, 1.8, 3.6, 5.0, "Plex / Jellyfin Silo", [
        "Private Island: Files trapped on local server; zero global or community discovery.",
        "NAT / CGNAT Barrier: Most residential ISPs block incoming ports; requires risky port forwarding.",
        "Self-Hosting Burden: Friends & family cannot find videos unless manually granted individual access.",
        "Overwhelms Home Upload: Every concurrent stream drains home upload bandwidth simultaneously."
    ], tag="TRADITIONAL SELF-HOSTING", tag_color=COLOR_RED, border_color=COLOR_RED)

    c2 = add_card(s2, 4.8, 1.8, 3.6, 5.0, "YouTube / Cloud Silo", [
        "Astronomical Storage Costs: Multi-gigabyte videos stored permanently in centralized cloud buckets.",
        "Privacy & Ownership Loss: Video binaries surrendered to centralized tech giants.",
        "Arbitrary Moderation: Instant takedowns and algorithmic de-monetization without recourse.",
        "Egress Bandwidth Tax: High recurring bandwidth bills that make independent community apps unviable."
    ], tag="CENTRALIZED CLOUD", tag_color=COLOR_RED, border_color=COLOR_RED)

    c3 = add_card(s2, 8.8, 1.8, 3.733, 5.0, "CustardFin Innovation", [
        "Best of Both Worlds: Video binaries remain 100% at home; metadata coordinates globally.",
        "Zero Cloud Video Storage: Cloud only stores ~3 KB JSON manifests in PostgreSQL.",
        "Outbound-Only Tunnels: FRP reverse tunnel punches through residential CGNAT with 0 inbound router ports.",
        "Origin Shield Cache: Home uploads each hot chunk ONCE; S3 cache serves hundreds of remote viewers."
    ], tag="THE CUSTARDFIN ANSWER", tag_color=COLOR_EMERALD, border_color=COLOR_EMERALD)

    inject_shape_animations(s2, [c1, c2, c3])

    # =========================================================================
    # SLIDE 3: THE GOLDEN INVARIANT (Zero Cloud Video Storage)
    # =========================================================================
    s3 = add_slide_base(transition="fade")
    add_header(s3, "Invariant 1: Zero Cloud Video Storage Footprint", "ARCHITECTURAL FOUNDATION")

    c1 = add_card(s3, 0.8, 1.8, 5.6, 5.0, "Traditional Cloud Architecture", [
        "Every video uploaded requires multi-gigabyte permanent cloud storage (S3 / GCS).",
        "1,000 community videos = ~4 TB storage cost permanently billed every month.",
        "Egress bandwidth costs compound as audience grows.",
        "Centralized point of failure and single target for data breaches.",
        "Owner loses physical custody of original uncompressed media assets."
    ], tag="CENTRALIZED INEFFICIENCY", tag_color=COLOR_RED)

    c2 = add_card(s3, 6.8, 1.8, 5.733, 5.0, "CustardFin Invariant 1", [
        "Video files strictly stay on owner's hard drive and edge SQLite database.",
        "Marking a video 'Public' syncs ONLY a ~3 KB JSON manifest to PostgreSQL:",
        "  • media_id, content_hash (SHA-256), duration_s, title",
        "  • 4s HLS renditions array (1080p, 720p, 480p playlists)",
        "Origin Shield Pull-Through: Cloud temporarily buffers active chunks with 24h LRU auto-purge.",
        "Near-zero cloud bill: A single $5/month cloud VM can coordinate discovery for 100,000+ videos globally!"
    ], tag="SOVEREIGN EDGE STORAGE", tag_color=COLOR_GOLD, border_color=COLOR_GOLD)

    inject_shape_animations(s3, [c1, c2])

    # =========================================================================
    # SLIDE 4: TWO-TIER ARCHITECTURE
    # =========================================================================
    s4 = add_slide_base(transition="push", push_dir="l")
    add_header(s4, "Two-Tier System: Home Edge Node & Cloud Control Plane", "SYSTEM ARCHITECTURE")

    c1 = add_card(s4, 0.8, 1.8, 5.6, 5.0, "Tier 1: Home Edge Node (Go Daemon)", [
        "Resource Friendly: Single static Go binary (<20 MB idle RAM footprint).",
        "Edge SQLite: Pure-Go driver (WAL mode, busy_timeout) storing paths, hashes & private watch history.",
        "Managed Transcoder: FFmpeg subprocess with automatic NVENC probe and CPU fallback.",
        "Pre-encoded 4s HLS: Chunks generated upfront to avoid live CPU exhaustion during playback.",
        "Outbound FRP Client: Connects outward to FRP server over TLS. Zero port forwarding required.",
        "Local Streaming Server: Serves direct LAN viewers on port 8090 at zero bandwidth cost."
    ], tag="EDGE NODE (HOME PC)", tag_color=COLOR_CYAN)

    c2 = add_card(s4, 6.8, 1.8, 5.733, 5.0, "Tier 2: Cloud Control Plane (Cloud VM / Lab)", [
        "Cloud API Daemon (Go): High throughput REST router and goroutine background workers.",
        "PostgreSQL (pgx + sqlc): Stores user accounts, node registry, and lightweight catalog manifests.",
        "Redis Heartbeat Cache: Stores ephemeral 'heartbeat:{node_id}' with 15s TTL for node online logic.",
        "FRP Relay Server: Bridges incoming remote playback requests to home edge tunnel.",
        "S3 Origin Shield: Pull-through chunk cache (MinIO / Floci) with signed HMAC URL tokens.",
        "Next.js 15 Web UI: Modern App Router frontend with hls.js player, search & Dev Console."
    ], tag="CLOUD CONTROL PLANE", tag_color=COLOR_GOLD)

    inject_shape_animations(s4, [c1, c2])

    # =========================================================================
    # SLIDE 5: WHAT IT USES (Technology Stack)
    # =========================================================================
    s5 = add_slide_base(transition="fade")
    add_header(s5, "Technology Stack: 100% Free & Open Source Dependencies", "WHAT IT USES")

    c1 = add_card(s5, 0.8, 1.8, 2.7, 5.0, "Backend / Edge", [
        "Go (latest stable):",
        "  • stdlib net/http",
        "  • log/slog structured logs",
        "  • context.Context everywhere",
        "SQLite (modernc.org):",
        "  • Pure Go (zero CGO)",
        "  • WAL mode & busy_timeout",
        "FFmpeg Managed Exec:",
        "  • NVENC hardware probe",
        "  • x264 CPU fallback",
        "  • Keyframes locked to 4.0s"
    ], tag="EDGE CORE", tag_color=COLOR_CYAN)

    c2 = add_card(s5, 3.8, 1.8, 2.7, 5.0, "Cloud Control", [
        "Go Cloud API (apid):",
        "  • Hexagonal architecture",
        "  • goroutines for async tasks",
        "PostgreSQL 16:",
        "  • pgx connection pool",
        "  • sqlc type-safe queries",
        "  • Goose versioned migrations",
        "Redis 7:",
        "  • go-redis client",
        "  • Heartbeats only (Invariant 5)"
    ], tag="CLOUD BACKEND", tag_color=COLOR_GOLD)

    c3 = add_card(s5, 6.8, 1.8, 2.7, 5.0, "Cache & Network", [
        "FRP (Fast Reverse Proxy):",
        "  • Outbound TLS tunnels",
        "  • Automatic backoff reconnect",
        "S3 Origin Shield:",
        "  • AWS SDK for Go v2",
        "  • MinIO or Floci emulator",
        "  • 24h LRU chunk caching",
        "Security & Auth:",
        "  • Argon2id password hashing",
        "  • HMAC-SHA256 signed URLs"
    ], tag="INFRASTRUCTURE", tag_color=COLOR_EMERALD)

    c4 = add_card(s5, 9.8, 1.8, 2.733, 5.0, "Frontend UI", [
        "Next.js 15 App Router:",
        "  • Server Components (SSR)",
        "  • TypeScript strict mode",
        "  • Node LTS runtime (no Bun)",
        "Video Player:",
        "  • hls.js adaptive client",
        "  • Live 'Stats for Nerds' HUD",
        "Design System:",
        "  • Pure CSS custom tokens",
        "  • Dark theme glassmorphism",
        "  • Fully responsive UI"
    ], tag="WEB CLIENT", tag_color=COLOR_GOLD)

    inject_shape_animations(s5, [c1, c2, c3, c4])

    # =========================================================================
    # SLIDE 6: WHAT IT DOES (Ingestion & 4s HLS Transcoding)
    # =========================================================================
    s6 = add_slide_base(transition="push", push_dir="l")
    add_header(s6, "What It Does: Video Ingestion & Hardware 4s HLS Pipeline", "CORE WORKFLOWS")

    c1 = add_card(s6, 0.8, 1.8, 2.7, 5.0, "1. Ingestion & Scan", [
        "Flexible Ingestion:",
        "  • Drag-and-drop web upload",
        "  • Direct folder drop (./media)",
        "  • 1-click test fixture generator",
        "Metadata Extraction:",
        "  • FFprobe codec inspection",
        "  • Width, height, duration",
        "  • Content hashing (SHA-256) to avoid re-encoding"
    ], tag="STEP 1", tag_color=COLOR_CYAN)

    c2 = add_card(s6, 3.8, 1.8, 2.7, 5.0, "2. Pre-Encoded HLS", [
        "Invariant 6 Compliance:",
        "  • Pre-encoded 4.0s chunks",
        "  • ZERO live on-the-fly transcode",
        "Multi-Rendition Ladder:",
        "  • 1080p (5000 kbps)",
        "  • 720p (2800 kbps)",
        "  • 480p (1200 kbps)",
        "Aligned Keyframes:",
        "  • -g 120 -keyint_min 120",
        "  • Seamless rendition switching"
    ], tag="STEP 2", tag_color=COLOR_GOLD)

    c3 = add_card(s6, 6.8, 1.8, 2.7, 5.0, "3. GPU Acceleration", [
        "Dual-Path Transcoder:",
        "  • Primary: NVIDIA NVENC (h264_nvenc)",
        "  • Fallback: CPU libx264",
        "Boot-Time Verification:",
        "  • 1s test-encode verification",
        "  • Checks hardware session limits",
        "Performance:",
        "  • ~12.5x real-time on RTX 4060",
        "  • Frees CPU for daemon tasks"
    ], tag="STEP 3", tag_color=COLOR_EMERALD)

    c4 = add_card(s6, 9.8, 1.8, 2.733, 5.0, "4. 3 KB Manifest Sync", [
        "Instant Publishing:",
        "  • 1-click Public toggle in UI",
        "  • Generates signed JSON manifest",
        "  • PUT /v1/nodes/{id}/media/{mid}",
        "Cloud Registration:",
        "  • Validated against JSON schema",
        "  • Upserted to PostgreSQL catalog",
        "  • Instantly discoverable by all community viewers"
    ], tag="STEP 4", tag_color=COLOR_GOLD)

    inject_shape_animations(s6, [c1, c2, c3, c4])

    # =========================================================================
    # SLIDE 7: INTELLIGENT PLAY DECISION ENGINE
    # =========================================================================
    s7 = add_slide_base(transition="fade")
    add_header(s7, "Intelligent Play Decision Engine (Stream Router)", "ROUTING & PLAYBACK DECISION")

    c1 = add_card(s7, 0.8, 1.8, 2.7, 5.0, "Step 1: Node Liveness", [
        "Checks Redis Heartbeat:",
        "  • Key: heartbeat:{node_id}",
        "  • TTL: 15 seconds",
        "If Key Missing:",
        "  • Node is marked offline",
        "  • Fallback: Serves cached chunks from Origin Shield",
        "  • If uncached: Returns clean 'Node Offline' notice"
    ], tag="STEP 1: LIVENESS", tag_color=COLOR_GOLD)

    c2 = add_card(s7, 3.8, 1.8, 2.7, 5.0, "Step 2: Subnet Detection", [
        "Compares IP Subnets:",
        "  • Checks viewer remote IP against node local CIDR",
        "If on Same Subnet (LAN):",
        "  • Routes viewer DIRECTLY to node's local HTTP streamer (:8090)",
        "  • 0 ms cloud relay latency",
        "  • ZERO cloud egress bandwidth",
        "  • Works even if internet drops!"
    ], tag="STEP 2: DIRECT LAN", tag_color=COLOR_EMERALD)

    c3 = add_card(s7, 6.8, 1.8, 2.7, 5.0, "Step 3: Cache Hit", [
        "Remote Viewer Request:",
        "  • Checks Origin Shield S3 cache for requested HLS chunk",
        "If Chunk Present (Cache HIT):",
        "  • Issues short-lived signed URL (HMAC-SHA256, 1h expiry)",
        "  • Served directly from S3 cache",
        "  • ZERO upload bytes consumed from the home owner's internet!"
    ], tag="STEP 3: ORIGIN SHIELD", tag_color=COLOR_CYAN)

    c4 = add_card(s7, 9.8, 1.8, 2.733, 5.0, "Step 4: Cache Miss", [
        "If Chunk Not In Cache:",
        "  • Checks admission control cap (max 10 active streams)",
        "  • Pulls chunk ONCE via FRP reverse tunnel from home node",
        "  • Writes chunk to S3 cache",
        "  • Issues signed playback URL",
        "Subsequent Viewers:",
        "  • Hit cache instantly (Step 3)"
    ], tag="STEP 4: PULL-THROUGH", tag_color=COLOR_GOLD)

    inject_shape_animations(s7, [c1, c2, c3, c4])

    # =========================================================================
    # SLIDE 8: SCALABILITY & ADMISSION CONTROL
    # =========================================================================
    s8 = add_slide_base(transition="push", push_dir="l")
    add_header(s8, "Scalability: Origin Shield & 10-Stream Concurrency Cap", "INVARIANTS 7 & 10")

    c1 = add_card(s8, 0.8, 1.8, 5.6, 5.0, "Origin Shield Invariant 7 (Pull-Through Cache)", [
        "The Viral Video Problem in Traditional P2P:",
        "  • If 50 people watch your video at once, home upload collapses.",
        "CustardFin Origin Shield Cache:",
        "  • Home edge node uploads each 4s video segment ONCE.",
        "  • Central S3-compatible cache stores chunk with 24h LRU policy.",
        "  • 100 viewers watching the same title = 99% cache hit ratio!",
        "  • Home ISP sees load equivalent to exactly 1 single viewer.",
        "  • Verified via load test scripts (scripts/load-10-viewers.ps1)."
    ], tag="INVARIANT 7: CACHE SHIELD", tag_color=COLOR_EMERALD, border_color=COLOR_EMERALD)

    c2 = add_card(s8, 6.8, 1.8, 5.733, 5.0, "Admission Control Invariant 10 (10-Stream Cap)", [
        "The 10 Distinct Titles Scenario:",
        "  • When 10 viewers watch 10 DIFFERENT titles, cache hit is 0%.",
        "  • Home upload speed divided by stream bitrate is the physical limit.",
        "CustardFin Admission Limiter:",
        "  • Strict per-node concurrency tracker in Go Cloud API.",
        "  • Hard limit: 10 concurrent distinct remote streams per node.",
        "  • 11th remote viewer receives: HTTP 429 Too Many Requests.",
        "  • Header: Retry-After: 30 ('Node Busy' clear message).",
        "  • Dynamic ABR ladder pruning: strips 1080p to keep 720p/480p alive."
    ], tag="INVARIANT 10: CAPACITY LIMIT", tag_color=COLOR_GOLD, border_color=COLOR_GOLD)

    inject_shape_animations(s8, [c1, c2])

    # =========================================================================
    # SLIDE 9: HARDWARE ACCELERATION & BENCHMARKS
    # =========================================================================
    s9 = add_slide_base(transition="fade")
    add_header(s9, "Hardware Acceleration: NVIDIA NVENC vs CPU Benchmark", "INVARIANT 6 & MILESTONE 9")

    c1 = add_card(s9, 0.8, 1.8, 5.6, 5.0, "Transcoding Engine Architecture", [
        "Probe Order:",
        "  1. NVIDIA NVENC (h264_nvenc) - Dedicated hardware ASIC chip.",
        "  2. Intel QSV / VAAPI (h264_qsv / h264_vaapi) - Integrated GPU.",
        "  3. Software CPU (libx264) - Automatic seamless fallback.",
        "Strict 4-Second Alignment:",
        "  • -g 120 -keyint_min 120 -sc_threshold 0",
        "  • Guarantees I-frame at start of every 4s HLS segment (.ts).",
        "Safe Fallback Guarantee:",
        "  • 1-second test probe on startup detects NVENC session limits.",
        "  • On any NVENC error, instantly falls back to CPU without crashing."
    ], tag="ARCHITECTURE", tag_color=COLOR_CYAN)

    c2 = add_card(s9, 6.8, 1.8, 5.733, 5.0, "Benchmark Results (RTX 4060 vs CPU)", [
        "Test Video: 15.0s 720p 30fps H.264 Test Fixture",
        "",
        "NVIDIA NVENC (Hardware ASIC):",
        "  • Transcode Time: 1.21 seconds (12.4x real-time speed)",
        "  • CPU Utilization: < 3% during encoding",
        "  • Chunk Duration: Exact 4.00s keyframe alignment",
        "",
        "Software libx264 (Multi-threaded CPU):",
        "  • Transcode Time: 5.42 seconds (2.7x real-time speed)",
        "  • CPU Utilization: 88% - 95% peak load across all cores",
        "",
        "Conclusion: NVENC delivers 4.5x faster encode speed and preserves CPU for edge streaming and OS responsiveness."
    ], tag="BENCHMARK DATA", tag_color=COLOR_GOLD, border_color=COLOR_GOLD)

    inject_shape_animations(s9, [c1, c2])

    # =========================================================================
    # SLIDE 10: MULTI-TOPOLOGY CLOUD DEPLOYMENT
    # =========================================================================
    s10 = add_slide_base(transition="push", push_dir="l")
    add_header(s10, "Deployment: Pure Local, Dev Containers, Floci & GCP", "DEPLOYMENT PROFILES")

    c1 = add_card(s10, 0.8, 1.8, 2.7, 5.0, "Mode 1: Pure Local", [
        "Setup:",
        "  • go run ./apps/api/cmd/apid",
        "  • go run ./apps/edge/cmd/edged",
        "  • npm run dev (Next.js)",
        "Use Case:",
        "  • Rapid developer iteration",
        "  • Zero external software required",
        "  • In-memory stores fallback"
    ], tag="MODE 1", tag_color=COLOR_CYAN)

    c2 = add_card(s10, 3.8, 1.8, 2.7, 5.0, "Mode 2: Dev Containers", [
        "Setup:",
        "  • docker compose -f deploy/compose.dev.yml up",
        "Services:",
        "  • PostgreSQL (port 5433 mapped)",
        "  • Redis 7 (port 6379)",
        "  • MinIO S3 (port 9000/9001)",
        "Use Case:",
        "  • Full local cloud simulation"
    ], tag="MODE 2", tag_color=COLOR_GOLD)

    c3 = add_card(s10, 6.8, 1.8, 2.7, 5.0, "Mode 3: Floci Emulator", [
        "Setup:",
        "  • docker compose -f deploy/floci/compose.floci.yml up",
        "Features:",
        "  • MIT open-source local cloud",
        "  • AWS S3 & GCP cloud emulation",
        "  • Multi-PC demonstration on lab Wi-Fi without internet"
    ], tag="MODE 3 (LAB DEMO)", tag_color=COLOR_EMERALD)

    c4 = add_card(s10, 9.8, 1.8, 2.733, 5.0, "Mode 4: Live GCP VM", [
        "Setup:",
        "  • Cloud VM (e2-standard-2)",
        "  • Caddy TLS reverse proxy",
        "  • frps on static public IP",
        "Use Case:",
        "  • Real-world global streaming",
        "  • Zero residential IP exposure",
        "  • HTTPS domain deployment"
    ], tag="MODE 4 (PRODUCTION)", tag_color=COLOR_GOLD)

    inject_shape_animations(s10, [c1, c2, c3, c4])

    # =========================================================================
    # SLIDE 11: LIVE DEMONSTRATION & REHEARSAL
    # =========================================================================
    s11 = add_slide_base(transition="fade")
    add_header(s11, "Live Demonstration Walkthrough & Script Scenarios", "DEMO REHEARSAL")

    c1 = add_card(s11, 0.8, 1.8, 5.6, 5.0, "5-Terminal Presentation Orchestration", [
        "Synchronized Presentation Launcher (scripts/start-all-terminals.ps1):",
        "  • Tab 1 [DOCKER]: compose.dev.yml (Postgres, Redis, MinIO)",
        "  • Tab 2 [CLOUD API]: apid listening on :8080 with auto-seeding",
        "  • Tab 3 [EDGE DAEMON]: edged with NVIDIA NVENC hardware transcoder",
        "  • Tab 4 [WEB UI]: Next.js 15 client on http://localhost:3000",
        "  • Tab 5 [SIMULATION]: Load test & benchmark script launcher",
        "Color-Coded Status Logs:",
        "  • Real-time heartbeat logs, manifest syncs, and chunk fetches."
    ], tag="ORCHESTRATION", tag_color=COLOR_CYAN)

    c2 = add_card(s11, 6.8, 1.8, 5.733, 5.0, "Core Demonstration Scenarios", [
        "Scenario 1: Direct LAN Playback (0ms overhead, 0 cloud bandwidth).",
        "Scenario 2: Remote Playback via FRP tunnel + Origin Shield cache.",
        "Scenario 3: 10 Viewers on Same Title: Origin Shield absorbs 90%+ upload load.",
        "Scenario 4: 10 Different Titles: Hits admission cap gracefully (HTTP 429).",
        "Scenario 5: Node Disconnect: Edge daemon killed -> catalog displays cached status.",
        "Scenario 6: Hardware NVENC vs CPU benchmark comparison.",
        "Interactive UI: Add Video modal, Edge Library manager, Dev Settings console & Stats for Nerds HUD."
    ], tag="VIVA SCENARIOS", tag_color=COLOR_GOLD, border_color=COLOR_GOLD)

    inject_shape_animations(s11, [c1, c2])

    # =========================================================================
    # SLIDE 12: CONCLUSION & SUMMARY
    # =========================================================================
    s12 = add_slide_base(transition="fade")
    add_header(s12, "Summary: Sovereign Media with Cloud-Scale Economics", "CONCLUSION & Q&A")

    c1 = add_card(s12, 0.8, 1.8, 3.6, 5.0, "Architecture Win", [
        "Zero Cloud Video Storage:",
        "  • Videos stay on owner's hard drive.",
        "  • ~3 KB manifest coordinates catalog.",
        "  • Near-zero recurring cloud cost.",
        "Outbound-Only Networking:",
        "  • No open router ports (ISP CGNAT).",
        "  • FRP TLS reverse tunnel.",
        "Origin Shield Pull-Through:",
        "  • Absorbs viral viewer traffic."
    ], tag="CORE INNOVATION", tag_color=COLOR_EMERALD)

    c2 = add_card(s12, 4.8, 1.8, 3.6, 5.0, "Engineering Quality", [
        "100% Free & Open Source:",
        "  • Zero proprietary dependencies.",
        "Robust Testing:",
        "  • Unit tests for every package.",
        "  • Adapter contract tests.",
        "  • 10-viewer load simulation scripts.",
        "Clean Hexagonal Design:",
        "  • Strict domain/ports/adapters.",
        "  • 10 Milestones complete."
    ], tag="RIGOROUS STANDARDS", tag_color=COLOR_CYAN)

    c3 = add_card(s12, 8.8, 1.8, 3.733, 5.0, "Evaluation Ready", [
        "Project Status: COMPLETE",
        "Team Members:",
        "  • Anish Ghadi (25P0320009)",
        "  • Anosh Chodankar (25P0320010)",
        "Repository:",
        "  • github.com/Anosh26/Custard-Fin",
        "Documentation:",
        "  • BRD v2, 4 Deployment Runbooks,",
        "  • Live Walkthroughs & Architecture ADRs.",
        "",
        "Thank You! Questions & Live Demo."
    ], tag="MCA LAB 2026", tag_color=COLOR_GOLD, border_color=COLOR_GOLD)

    inject_shape_animations(s12, [c1, c2, c3])

    # Save Presentation
    output_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "docs", "presentation"))
    os.makedirs(output_dir, exist_ok=True)
    pptx_path = os.path.join(output_dir, "CustardFin_Presentation.pptx")
    prs.save(pptx_path)
    print(f"[OK] Generated PowerPoint presentation: {pptx_path}")

    # Export OpenXML Slide XMLs
    xml_path = os.path.join(output_dir, "presentation.xml")
    export_presentation_xml(pptx_path, xml_path)
    print(f"[OK] Exported OpenXML presentation XML: {xml_path}")

    return pptx_path, xml_path

def export_presentation_xml(pptx_path, xml_out_path):
    """Extracts and formats the OpenXML presentation structure into a readable consolidated XML file."""
    root = ET.Element("PresentationML", {
        "title": "CustardFin Presentation",
        "aspectRatio": "16:9",
        "theme": "Dark Modern (Gold/Cyan/Emerald)",
        "xmlns:p": "http://schemas.openxmlformats.org/presentationml/2006/main",
        "xmlns:a": "http://schemas.openxmlformats.org/drawingml/2006/main"
    })

    with zipfile.ZipFile(pptx_path, 'r') as z:
        # Read presentation metadata
        if 'ppt/presentation.xml' in z.namelist():
            pres_xml = z.read('ppt/presentation.xml')
            p_elem = ET.fromstring(pres_xml)
            meta = ET.SubElement(root, "PresentationProperties")
            meta.append(p_elem)

        # Read all slides with their timing animations and transitions
        slide_files = sorted([f for f in z.namelist() if f.startswith('ppt/slides/slide') and f.endswith('.xml')],
                             key=lambda x: int(''.join(filter(str.isdigit, x)) or 0))
        
        slides_elem = ET.SubElement(root, "Slides", {"count": str(len(slide_files))})

        for sf in slide_files:
            slide_data = z.read(sf)
            s_elem = ET.fromstring(slide_data)
            slide_node = ET.SubElement(slides_elem, "Slide", {"filename": sf})
            slide_node.append(s_elem)

    # Format pretty XML
    rough_str = ET.tostring(root, 'utf-8')
    reparsed = minidom.parseString(rough_str)
    pretty_xml = reparsed.toprettyxml(indent="  ")

    with open(xml_out_path, "w", encoding="utf-8") as f:
        f.write(pretty_xml)

if __name__ == "__main__":
    create_presentation()
