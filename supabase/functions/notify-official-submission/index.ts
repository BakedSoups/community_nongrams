import { serve } from "https://deno.land/std@0.224.0/http/server.ts";

serve(async (request) => {
  const expectedSecret = Deno.env.get("WEBHOOK_SECRET");
  if (!expectedSecret || request.headers.get("x-webhook-secret") !== expectedSecret) {
    return new Response("unauthorized", { status: 401 });
  }
  const payload = await request.json();
  const submission = payload.record;
  const resendKey = Deno.env.get("RESEND_API_KEY");
  const reviewEmail = Deno.env.get("REVIEW_EMAIL");
  const adminURL = Deno.env.get("ADMIN_REVIEW_URL");
  if (!resendKey || !reviewEmail || !adminURL || !submission?.id) {
    return new Response("missing function configuration", { status: 500 });
  }
  const response = await fetch("https://api.resend.com/emails", {
    method: "POST",
    headers: { Authorization: `Bearer ${resendKey}`, "Content-Type": "application/json" },
    body: JSON.stringify({
      from: Deno.env.get("REVIEW_FROM_EMAIL") || "Pixaross <submissions@example.com>",
      to: [reviewEmail],
      subject: "New Pixaross main-game submission",
      html: `<p>A creator submitted a level for inspection.</p><p><a href="${adminURL}?submission=${encodeURIComponent(submission.id)}">Review submission</a></p>`
    })
  });
  return new Response(await response.text(), { status: response.status, headers: { "content-type": "application/json" } });
});
