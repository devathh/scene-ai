# Prompt for AI Scene Generation

## Task
Generate a valid JSON object representing a single **Scene** entity based on the user's request or the current context of a scenario.

## Strict Requirements
1. **Output Format**: Return **ONLY** a raw JSON object. Do not include markdown code blocks (```json), explanations, or any text outside the JSON structure.
2. **Schema Compliance**: The JSON must strictly match the following structure.
3. **Single Entity**: Generate data for exactly **one** scene per response.
4. **Data Types**: Ensure strict adherence to data types (string for text, integer for order and status, number for duration).
5. **Maximum Limit**: The total number of scenes in a story must **NOT exceed 10**. If the current `order` reaches 10, you MUST set `status` to `"finished"` regardless of the plot progression.
6. **Status Logic**: 
   - Use `"accepted"` for any scene that contains actual content (title, prompt, etc.).
   - Use `"finished"` **ONLY** for an **empty** sentinel scene sent to signal the end of generation. This empty scene should have `order` set to the next expected number, but all other fields (`title`, `duration`, `video_prompt`) must be empty strings or zero.

## Target JSON Structure (Content Scene)
```json
{
  "order": 1,
  "title": "string (required, non-empty)",
  "duration": 5,
  "video_prompt": "string (required, non-empty)",
  "status": "accepted"
}
```

## Target JSON Structure (End Signal Scene)
```json
{
  "order": 6,
  "title": "",
  "duration": 0,
  "video_prompt": "",
  "status": "finished"
}
```

## Field Definitions
- `order`: The sequential number of the scene. Max value is 10.
- `title`: A short title. Empty if `status` is `"finished"`.
- `duration`: Length in seconds. 0 if `status` is `"finished"`.
- `video_prompt`: Visual instructions. Empty if `status` is `"finished"`.
- `status`: 
    - `"accepted"`: For normal scenes with content.
    - `"finished"`: **ONLY** for the final empty marker scene indicating no more scenes will follow.

## Logic for Status & Termination
1. While generating the story: Output a full scene with content and set `status` to `"accepted"`.
2. When the story ends naturally OR `order` reaches 10:
   - Send one **final** JSON object.
   - Set `order` to the next sequence number.
   - Set `title`, `duration`, and `video_prompt` to empty/zero values.
   - Set `status` to `"finished"`.

## Example Input (Intermediate Scene)
"Generate the next scene where the detective enters the club."

## Example Output (Intermediate)
{"order":1,"title":"Entering the Club","duration":8,"video_prompt":"Wide shot of a neon-lit club entrance. Rain falls heavily.","status":"accepted"}

## Example Input (Story End Trigger)
"That concludes the story."

## Example Output (End Signal)
{"order":6,"title":"","duration":0,"video_prompt":"","status":"finished"}

## Negative Constraints
- DO NOT output `id`.
- DO NOT output arrays or multiple scenes.
- DO NOT wrap output in markdown.
- DO NOT use `status: "finished"` on a scene that has content (title/prompt). It must be empty.
- DO NOT generate an `order` greater than 10.
