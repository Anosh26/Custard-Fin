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

COLOR_BG = RGBColor(7, 10, 18)
COLOR_BG_CARD = RGBColor(17, 24, 39)
COLOR_BORDER_DEFAULT = RGBColor(51, 65, 85)
COLOR_TEXT_MAIN = RGBColor(248, 250, 252)
COLOR_TEXT_MUTED = RGBColor(203, 213, 225)
COLOR_TEXT_DIM = RGBColor(148, 163, 184)

COLOR_PASTEL_CUSTARD = RGBColor(254, 240, 138)
COLOR_PASTEL_CARAMEL = RGBColor(253, 224, 71)
COLOR_PASTEL_MINT = RGBColor(167, 243, 208)
COLOR_PASTEL_LAVENDER = RGBColor(221, 214, 254)
COLOR_PASTEL_SKY = RGBColor(186, 230, 253)
COLOR_PASTEL_PEACH = RGBColor(254, 215, 170)
COLOR_PASTEL_ROSE = RGBColor(254, 205, 211)

FONT_HEADING = "Plus Jakarta Sans"
FONT_BODY = "Plus Jakarta Sans"

def create_presentation():
    prs = Presentation()
    prs.slide_width = Inches(13.333)
    prs.slide_height = Inches(7.5)
    blank_layout = prs.slide_layouts[6]

    base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
    mascot_png = os.path.join(base_dir, "docs", "presentation", "custard_mascot.png")
    has_mascot = os.path.exists(mascot_png)

    def add_slide_base(transition="fade", push_dir="l", mascot_speech=""):
        slide = prs.slides.add_slide(blank_layout)

        if transition == "push":
            trans_xml = f'<p:transition {nsdecls("p")} spd="med"><p:push dir="{push_dir}"/></p:transition>'
        elif transition == "wipe":
            trans_xml = f'<p:transition {nsdecls("p")} spd="med"><p:wipe dir="r"/></p:transition>'
        else:
            trans_xml = f'<p:transition {nsdecls("p")} spd="med"><p:fade filter="out"/></p:transition>'
        slide.element.append(parse_xml(trans_xml))

        bg = slide.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, 0, Inches(13.333), Inches(7.5))
        bg.fill.solid()
        bg.fill.fore_color.rgb = COLOR_BG
        bg.line.fill.background()

        top_bar = slide.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, 0, Inches(13.333), Inches(0.08))
        top_bar.fill.solid()
        top_bar.fill.fore_color.rgb = COLOR_PASTEL_CUSTARD
        top_bar.line.fill.background()

        if mascot_speech:
            bubble = slide.shapes.add_shape(MSO_SHAPE.ROUNDED_RECTANGLE, Inches(9.2), Inches(6.52), Inches(2.75), Inches(0.52))
            bubble.fill.solid()
            bubble.fill.fore_color.rgb = RGBColor(19, 27, 46)
            bubble.line.color.rgb = COLOR_PASTEL_CUSTARD
            bubble.line.width = Pt(1.2)
            tf_b = bubble.text_frame
            tf_b.word_wrap = False
            p_b = tf_b.paragraphs[0]
            p_b.alignment = PP_ALIGN.CENTER
            p_b.text = mascot_speech
            p_b.font.name = FONT_HEADING
            p_b.font.size = Pt(10.5)
            p_b.font.bold = True
            p_b.font.color.rgb = COLOR_PASTEL_CUSTARD

        if has_mascot:
            slide.shapes.add_picture(mascot_png, Inches(12.1), Inches(6.18), width=Inches(0.98), height=Inches(0.98))

        return slide

    def add_header(slide, title, category="CUSTARDFIN • DISTRIBUTED COMMUNITY MEDIA", category_color=COLOR_PASTEL_CUSTARD):
        tb_cat = slide.shapes.add_textbox(Inches(0.8), Inches(0.35), Inches(11.0), Inches(0.32))
        tf_cat = tb_cat.text_frame
        tf_cat.word_wrap = True
        p_cat = tf_cat.paragraphs[0]
        p_cat.text = category.upper()
        p_cat.font.name = FONT_HEADING
        p_cat.font.size = Pt(11)
        p_cat.font.bold = True
        p_cat.font.color.rgb = category_color

        tb_title = slide.shapes.add_textbox(Inches(0.8), Inches(0.65), Inches(11.0), Inches(0.75))
        tf_title = tb_title.text_frame
        tf_title.word_wrap = True
        p_title = tf_title.paragraphs[0]
        p_title.text = title
        p_title.font.name = FONT_HEADING
        p_title.font.size = Pt(27)
        p_title.font.bold = True
        p_title.font.color.rgb = COLOR_TEXT_MAIN

        divider = slide.shapes.add_shape(MSO_SHAPE.RECTANGLE, Inches(0.8), Inches(1.45), Inches(11.733), Inches(0.02))
        divider.fill.solid()
        divider.fill.fore_color.rgb = RGBColor(40, 53, 76)
        divider.line.fill.background()

    def add_card(slide, x, y, w, h, title, items, tag=None, tag_color=COLOR_PASTEL_CUSTARD, border_color=None, metric_val=None, metric_lbl=None):
        if border_color is None:
            border_color = tag_color

        card = slide.shapes.add_shape(MSO_SHAPE.ROUNDED_RECTANGLE, Inches(x), Inches(y), Inches(w), Inches(h))
        card.fill.solid()
        card.fill.fore_color.rgb = COLOR_BG_CARD
        card.line.color.rgb = border_color
        card.line.width = Pt(1.5)

        pad_x = 0.28
        pad_y = 0.22
        tb = slide.shapes.add_textbox(Inches(x + pad_x), Inches(y + pad_y), Inches(w - 2 * pad_x), Inches(h - 2 * pad_y))
        tf = tb.text_frame
        tf.word_wrap = True

        if tag:
            p_tag = tf.paragraphs[0]
            p_tag.text = tag.upper()
            p_tag.font.name = FONT_HEADING
            p_tag.font.size = Pt(10)
            p_tag.font.bold = True
            p_tag.font.color.rgb = tag_color
            p_tag.space_after = Pt(4)
            p_title = tf.add_paragraph()
        else:
            p_title = tf.paragraphs[0]

        p_title.text = title
        p_title.font.name = FONT_HEADING
        p_title.font.size = Pt(18)
        p_title.font.bold = True
        p_title.font.color.rgb = COLOR_TEXT_MAIN
        p_title.space_after = Pt(8)

        if metric_val:
            p_mval = tf.add_paragraph()
            p_mval.text = metric_val
            p_mval.font.name = FONT_HEADING
            p_mval.font.size = Pt(26)
            p_mval.font.bold = True
            p_mval.font.color.rgb = tag_color
            p_mval.space_after = Pt(2)
            if metric_lbl:
                p_mlbl = tf.add_paragraph()
                p_mlbl.text = metric_lbl
                p_mlbl.font.name = FONT_BODY
                p_mlbl.font.size = Pt(11)
                p_mlbl.font.color.rgb = COLOR_TEXT_DIM
                p_mlbl.space_after = Pt(8)

        for item in items:
            p = tf.add_paragraph()
            p.text = f"▸  {item}"
            p.font.name = FONT_BODY
            p.font.size = Pt(13.5)
            p.font.color.rgb = COLOR_TEXT_MUTED
            p.space_after = Pt(7)

        return card.shape_id

    def add_progress_bar(slide, x, y, w, h, pct, fill_color):
        track = slide.shapes.add_shape(MSO_SHAPE.ROUNDED_RECTANGLE, Inches(x), Inches(y), Inches(w), Inches(h))
        track.fill.solid()
        track.fill.fore_color.rgb = RGBColor(30, 41, 59)
        track.line.fill.background()

        bar_w = max(0.1, w * min(1.0, max(0.0, pct)))
        bar = slide.shapes.add_shape(MSO_SHAPE.ROUNDED_RECTANGLE, Inches(x), Inches(y), Inches(bar_w), Inches(h))
        bar.fill.solid()
        bar.fill.fore_color.rgb = fill_color
        bar.line.fill.background()

    def inject_shape_animations(slide, shape_ids):
        if not shape_ids:
            return

        child_nodes = []
        c_id = 3
        for i, spid in enumerate(shape_ids):
            node_type = "clickEffect" if i == 0 else "afterEffect"
            delay = 0 if i == 0 else 220
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

    s1 = add_slide_base(transition="fade", mascot_speech="Welcome to CustardFin! ✨")
    tb1 = s1.shapes.add_textbox(Inches(0.8), Inches(0.9), Inches(8.8), Inches(3.0))
    tf1 = tb1.text_frame
    tf1.word_wrap = True

    p_pill = tf1.paragraphs[0]
    p_pill.text = "MCA CLOUD COMPUTING LAB PROJECT • SEMESTER END EVALUATION"
    p_pill.font.name = FONT_HEADING
    p_pill.font.size = Pt(11)
    p_pill.font.bold = True
    p_pill.font.color.rgb = COLOR_PASTEL_CUSTARD
    p_pill.space_after = Pt(8)

    p_title = tf1.add_paragraph()
    p_title.text = "CUSTARDFIN: Community Video Platform"
    p_title.font.name = FONT_HEADING
    p_title.font.size = Pt(36)
    p_title.font.bold = True
    p_title.font.color.rgb = COLOR_PASTEL_CUSTARD

    p_sub = tf1.add_paragraph()
    p_sub.text = "Keep Videos at Home, Share with the World"
    p_sub.font.name = FONT_HEADING
    p_sub.font.size = Pt(21)
    p_sub.font.bold = True
    p_sub.font.color.rgb = COLOR_TEXT_MAIN
    p_sub.space_before = Pt(4)

    p_desc = tf1.add_paragraph()
    p_desc.text = "Keep your videos at home like Jellyfin, share them with the world like YouTube — with Zero Permanent Cloud Video Storage."
    p_desc.font.name = FONT_BODY
    p_desc.font.size = Pt(14)
    p_desc.font.color.rgb = COLOR_TEXT_MUTED
    p_desc.space_before = Pt(8)

    if has_mascot:
        s1.shapes.add_picture(mascot_png, Inches(10.1), Inches(0.85), width=Inches(2.4), height=Inches(2.4))

    c1 = add_card(s1, 0.8, 4.15, 3.65, 2.25, "Team Credits", [
        "Anish Ghadi (25P0320009)",
        "Anosh Chodankar (25P0320010)",
        "MCA Cloud Computing Lab 2026",
        "Self-Hosted Cloud Demonstration"
    ], tag="AUTHORS", tag_color=COLOR_PASTEL_LAVENDER)

    c2 = add_card(s1, 4.8, 4.15, 3.65, 2.25, "Core Rule", [
        "Videos stay 100% on home PC",
        "Public sharing syncs ~3 KB text file",
        "Zero permanent cloud storage bills",
        "Private watch history stays on edge"
    ], tag="INVARIANT #1", tag_color=COLOR_PASTEL_CUSTARD)

    c3 = add_card(s1, 8.8, 4.15, 3.733, 2.25, "Simple Tech Stack", [
        "Tiny Go 1.24 static binary (< 20 MB RAM)",
        "Next.js 15 App Router web client",
        "PostgreSQL + SQLite WAL + Redis",
        "Universal GPU + CPU video encoding"
    ], tag="100% OPEN SOURCE", tag_color=COLOR_PASTEL_MINT)

    inject_shape_animations(s1, [c1, c2, c3])

    s2 = add_slide_base(transition="push", push_dir="l", mascot_speech="Best of Both Worlds! 💡")
    add_header(s2, "Why Current Apps Fail: Big Cloud vs Home Servers", "THE CORE PROBLEM", COLOR_PASTEL_ROSE)

    c1 = add_card(s2, 0.8, 1.65, 5.65, 3.75, "Big Cloud (YouTube / AWS S3)", [
        "Monthly storage bills grow continuously as library expands.",
        "Creators surrender video ownership to corporate platform ads.",
        "Crushing cloud bandwidth fees make community apps expensive."
    ], tag="TRADITIONAL CLOUD", tag_color=COLOR_PASTEL_ROSE, metric_val="HEAVY FEES", metric_lbl="Compounding Monthly Costs")

    c2 = add_card(s2, 6.85, 1.65, 5.65, 3.75, "Home Servers (Plex / Jellyfin)", [
        "Blocked by ISP CGNAT; needs complicated port forwarding.",
        "Private silos with no public feed for friends to discover media.",
        "3 to 4 outside viewers choke home Wi-Fi upload bandwidth."
    ], tag="HOME SILO LIMITS", tag_color=COLOR_PASTEL_LAVENDER, metric_val="BLOCKED", metric_lbl="ISP CGNAT & Choked Upload")

    sol = add_card(s2, 0.8, 5.55, 11.7, 0.85, "The CustardFin Fix", [
        "Videos stay at home • Cloud only stores a ~3 KB info file • Cloud cache protects home Wi-Fi."
    ], tag="BEST OF BOTH WORLDS", tag_color=COLOR_PASTEL_CUSTARD)

    inject_shape_animations(s2, [c1, c2, sol])

    s3 = add_slide_base(transition="push", push_dir="l", mascot_speech="4 Simple Steps! ⚙️")
    add_header(s3, "How It Works: From Home Folder to Public Stream", "4-STEP PIPELINE", COLOR_PASTEL_CUSTARD)

    c1 = add_card(s3, 0.8, 1.65, 2.72, 4.75, "Scan & Hash", [
        "Scans local folder for video files.",
        "Generates thumbnail preview frame.",
        "Saves unique fingerprint in SQLite."
    ], tag="STEP 01 • HOME", tag_color=COLOR_PASTEL_CUSTARD, metric_val="SHA-256", metric_lbl="Unique Content ID")

    c2 = add_card(s3, 3.75, 1.65, 2.72, 4.75, "4s Video Chunks", [
        "Pre-encodes video into 4-second .ts parts.",
        "Hardware GPU or CPU fallback.",
        "Zero playback lag when users watch."
    ], tag="STEP 02 • ENCODE", tag_color=COLOR_PASTEL_MINT, metric_val="4s HLS", metric_lbl="Pre-Encoded Chunks")

    c3 = add_card(s3, 6.7, 1.65, 2.72, 4.75, "Safe Outbound Tunnel", [
        "Outbound FRP link connects to GCP VM.",
        "Bypasses CGNAT with no open ports.",
        "15s health ping proves node is online."
    ], tag="STEP 03 • CONNECT", tag_color=COLOR_PASTEL_SKY, metric_val="OUTBOUND", metric_lbl="Bypasses ISP CGNAT")

    c4 = add_card(s3, 9.65, 1.65, 2.85, 4.75, "3 KB Info Sync", [
        "1-click Public syncs JSON manifest.",
        "Saves title & duration to Postgres.",
        "Surfaces in community discovery feed."
    ], tag="STEP 04 • SHARE", tag_color=COLOR_PASTEL_PEACH, metric_val="~3 KB", metric_lbl="Manifest Payload")

    inject_shape_animations(s3, [c1, c2, c3, c4])

    s4 = add_slide_base(transition="fade", mascot_speech="Home & Cloud Together! 🤝")
    add_header(s4, "Two-Part Design: Home PC + Cloud Server", "SYSTEM ARCHITECTURE", COLOR_PASTEL_SKY)

    c1 = add_card(s4, 0.8, 1.65, 5.65, 4.75, "Part 1 • Home Edge App (Any PC)", [
        "Tiny Go app runs quietly in background using under 20 MB RAM.",
        "Local SQLite WAL DB stores paths, hashes, and private watch progress.",
        "FFmpeg worker creates 4s chunks and video preview thumbnails.",
        "FRP client connects outbound; LAN streamer serves home gigabit fast."
    ], tag="RUNS AT HOME", tag_color=COLOR_PASTEL_CUSTARD, metric_val="< 20 MB", metric_lbl="Idle RAM Footprint")

    c2 = add_card(s4, 6.85, 1.65, 5.65, 4.75, "Part 2 • Cloud Server (Google Cloud)", [
        "Next.js 15 web UI with video player and community feed.",
        "Go Cloud API handles logins and play decision routing.",
        "PostgreSQL holds ~3 KB manifests; Redis tracks 15s node heartbeats.",
        "24h Origin Shield chunk cache serves hot video pieces to viewers."
    ], tag="RUNS ON GCP", tag_color=COLOR_PASTEL_SKY, metric_val="0 VIDEOS", metric_lbl="Zero Permanent Cloud Storage")

    inject_shape_animations(s4, [c1, c2])

    s5 = add_slide_base(transition="push", push_dir="l", mascot_speech="Powered by GCP! ☁️")
    add_header(s5, "How We Use the Cloud: GCP Features & Setup", "GCP SETUP", COLOR_PASTEL_CUSTARD)

    c1 = add_card(s5, 0.8, 1.65, 2.72, 4.75, "GCP Compute Engine", [
        "Linux host running Docker containers.",
        "Runs Caddy, Next.js, API, and FRPS.",
        "Static public IP for reliable tunnel."
    ], tag="CLOUD COMPUTE", tag_color=COLOR_PASTEL_CUSTARD, metric_val="GCP VM", metric_lbl="Linux Container Host")

    c2 = add_card(s5, 3.75, 1.65, 2.72, 4.75, "VPC & Firewall", [
        "Ports 80 & 443 open for web traffic.",
        "Port 7000 open for safe FRP tunnel.",
        "Caddy provides automated HTTPS TLS."
    ], tag="NETWORKING", tag_color=COLOR_PASTEL_SKY, metric_val="VPC TLS", metric_lbl="Firewall & Auto HTTPS")

    c3 = add_card(s5, 6.7, 1.65, 2.72, 4.75, "In-VM S3 Object Cache", [
        "Runs inside the GCP VM via MinIO container.",
        "Exposes standard S3 API on port 9000.",
        "No separate paid cloud bucket needed.",
        "Auto-purges inactive chunks after 24 hours."
    ], tag="ORIGIN SHIELD", tag_color=COLOR_PASTEL_MINT, metric_val="24h LRU", metric_lbl="In-VM S3 Cache (0 Extra Cost)")

    c4 = add_card(s5, 9.65, 1.65, 2.85, 4.75, "Cloud State Layer", [
        "PostgreSQL stores ~3 KB manifests & users.",
        "Redis heartbeat keys have 15s TTL.",
        "Instant detection if home PC goes offline."
    ], tag="METADATA & STATE", tag_color=COLOR_PASTEL_PEACH, metric_val="15s TTL", metric_lbl="Postgres & Redis Health")

    inject_shape_animations(s5, [c1, c2, c3, c4])

    s6 = add_slide_base(transition="fade", mascot_speech="Home Wi-Fi Safe! 🛡️")
    add_header(s6, "Origin Shield & Router: Protecting Home Wi-Fi", "SMART ROUTING & SHIELD", COLOR_PASTEL_MINT)

    c1 = add_card(s6, 0.8, 1.65, 5.65, 4.75, "Smart Play Router (4-Step Check)", [
        "1. Redis Liveness: Checks if home node is online within last 15 seconds.",
        "2. Same Wi-Fi Bypass: Home viewers stream direct via LAN (0 cloud data).",
        "3. Cloud Cache Hit: Serves cached 4s chunk from GCP bucket instantly.",
        "4. Cache Miss Pull: Pulls chunk once through tunnel; caps at 10 streams."
    ], tag="PLAY DECISION ENGINE", tag_color=COLOR_PASTEL_CUSTARD, metric_val="4 CHECKS", metric_lbl="Intelligent Stream Routing")

    c2 = add_card(s6, 6.85, 1.65, 5.65, 4.75, "Home Wi-Fi Protection", [
        "Home PC uploads each hot 4s chunk once (5 Mbps) over the tunnel.",
        "GCP Origin Shield cache serves 100+ viral viewers directly.",
        "Home upload never chokes even when a video goes viral.",
        "Tested under 10 concurrent remote viewers with zero buffering."
    ], tag="BANDWIDTH SHIELD", tag_color=COLOR_PASTEL_MINT, metric_val="92.4% SAVED", metric_lbl="Home Upload Bandwidth Protected")

    add_progress_bar(s6, 7.15, 5.75, 5.0, 0.14, 0.924, COLOR_PASTEL_MINT)

    inject_shape_animations(s6, [c1, c2])

    s7 = add_slide_base(transition="push", push_dir="l", mascot_speech="Runs on Any PC! 🚀")
    add_header(s7, "Open-Source Tools & Universal Hardware Support", "UNIVERSAL HARDWARE", COLOR_PASTEL_CARAMEL)

    c1 = add_card(s7, 0.8, 1.65, 3.65, 4.75, "Smart Encoding", [
        "Auto-detects hardware on startup.",
        "Uses NVIDIA NVENC GPU if present.",
        "Automatic CPU (libx264) fallback.",
        "Runs on any team laptop or home PC."
    ], tag="UNIVERSAL HARDWARE", tag_color=COLOR_PASTEL_CUSTARD, metric_val="GPU / CPU", metric_lbl="NVENC + CPU libx264 Fallback")

    add_progress_bar(s7, 1.08, 5.75, 3.05, 0.14, 0.88, COLOR_PASTEL_CUSTARD)

    c2 = add_card(s7, 4.8, 1.65, 3.65, 4.75, "Low Memory", [
        "Static Go binary uses under 18 MB idle RAM.",
        "95% less RAM than Plex/Jellyfin (400+ MB).",
        "Embedded SQLite WAL for fast local storage.",
        "No fan noise or heat on team laptops."
    ], tag="RESOURCE EFFICIENCY", tag_color=COLOR_PASTEL_MINT, metric_val="< 18 MB", metric_lbl="Edge Daemon Idle RAM")

    add_progress_bar(s7, 5.08, 5.75, 3.05, 0.14, 0.15, COLOR_PASTEL_MINT)

    c3 = add_card(s7, 8.8, 1.65, 3.733, 4.75, "100% Free & Open Source", [
        "Go 1.24 static binary backend.",
        "Next.js 15 UI with React 19 & hls.js.",
        "PostgreSQL + SQLite WAL + Redis.",
        "FFmpeg, FRP tunnel, and Caddy TLS."
    ], tag="ZERO LICENSES", tag_color=COLOR_PASTEL_SKY, metric_val="100% FOSS", metric_lbl="Free & Open Source")

    inject_shape_animations(s7, [c1, c2, c3])

    s8 = add_slide_base(transition="fade", mascot_speech="Thank You! Questions? 🏆")
    add_header(s8, "Live Demonstration: What We Show Today", "VIVA READINESS", COLOR_PASTEL_CUSTARD)

    c1 = add_card(s8, 0.8, 1.65, 5.65, 4.75, "Live Demo Checklist", [
        "User signup and login with uploader attribution.",
        "Private vs Public toggle with instant ~3 KB manifest sync.",
        "Auto-generated thumbnails and self-cleaning library on boot.",
        "Direct LAN playback bypass (0 cloud data) + remote stream.",
        "Instant offline detection when home node shuts down (15s TTL)."
    ], tag="TESTED LIVE FEATURES", tag_color=COLOR_PASTEL_CUSTARD, metric_val="6 DEMOS", metric_lbl="End-to-End Live Scenarios")

    c2 = add_card(s8, 6.85, 1.65, 5.65, 4.75, "Project Summary", [
        "Zero permanent cloud video storage strictly preserved.",
        "Google Cloud Platform control plane manages discovery & cache.",
        "Universal hardware support runs on any team member's PC.",
        "Viva presentation defense: Anish Ghadi & Anosh Chodankar."
    ], tag="EVALUATION READY", tag_color=COLOR_PASTEL_MINT, metric_val="READY", metric_lbl="Viva Defense & Demo")

    inject_shape_animations(s8, [c1, c2])

    out_dir = os.path.join(base_dir, "docs", "presentation")
    os.makedirs(out_dir, exist_ok=True)
    pptx_path = os.path.join(out_dir, "CustardFin_Presentation.pptx")
    prs.save(pptx_path)
    print(f"Presentation saved successfully to: {pptx_path}")

    try:
        with zipfile.ZipFile(pptx_path, 'r') as z:
            xml_bytes = z.read("ppt/presentation.xml")
            parsed = minidom.parseString(xml_bytes)
            xml_pretty = parsed.toprettyxml(indent="  ")
            xml_path = os.path.join(out_dir, "presentation.xml")
            with open(xml_path, "w", encoding="utf-8") as f:
                f.write(xml_pretty)
            print(f"Presentation XML exported to: {xml_path}")
    except Exception as e:
        print(f"Warning: could not export raw XML: {e}")

if __name__ == "__main__":
    create_presentation()
