# Write Article

You are tasked with writing comprehensive article content based on research findings and a detailed plan, creating engaging and well-structured articles that serve the target audience effectively.

## CRITICAL: YOUR ROLE IS TO WRITE THE ACTUAL ARTICLE CONTENT
- DO NOT conduct additional research (use existing research and plan documents)
- DO NOT restructure the planned outline (follow the plan)
- DO NOT review or edit for final quality (that's for the review phase)
- FOCUS on writing clear, engaging, and informative content
- You are creating the first complete draft following the established plan

## Initial Setup:

When this command is invoked, respond with:
```
I'm ready to write your article. Please provide:
1. Path to the article plan document (from plan_article phase)
2. Path to research document (if available)
3. Target output location for the article draft
4. Any specific writing preferences or constraints
5. Whether to write the full article or specific sections
```

Then wait for the user's writing requirements.

## Steps to follow after receiving the writing requirements:

1. **Read plan and research documents:**
   - Read the article plan document FULLY first
   - **IMPORTANT**: Use the Read tool WITHOUT limit/offset parameters
   - Read the research document to understand available sources and findings
   - **CRITICAL**: Read all referenced materials before beginning writing
   - If plan document is missing, recommend running `/plan_article` first

2. **Analyze plan and prepare for writing:**
   - Review the article structure and section breakdown
   - Understand target audience, tone, and style requirements
   - Note key messages and content guidance for each section
   - Identify research integration points and source materials
   - Create writing strategy using TodoWrite to track section progress

3. **Write article content systematically:**
   - Follow the planned structure and section order
   - Write each section according to plan specifications
   - Integrate research findings, quotes, and examples as outlined
   - Maintain consistent tone and style throughout
   - Ensure smooth transitions between sections
   - Track progress and maintain focus on section objectives

4. **Gather metadata for the article document:**
   - Run the `hack/spec_metadata.sh` script to generate relevant metadata
   - Use filename as specified by user or follow pattern:
     - `articles/YYYY-MM-DD-topic-article.md` for thoughts directory
     - Or user-specified location for direct article output

5. **Generate complete article:**
   - Create article following plan structure with proper formatting
   - Include YAML frontmatter if saving to thoughts directory
   - Structure content for readability and engagement:
     ```markdown
     ---
     date: [Current date and time with timezone in ISO format]
     author: [Author name from thoughts status]
     git_commit: [Current commit hash]
     branch: [Current branch name]
     repository: [Repository name]
     title: "[Article Title]"
     topic: "[Article Topic]"
     target_audience: "[Target audience]"
     article_format: "[Format type]"
     word_count: [Actual word count]
     plan_document: "[Path to plan document]"
     research_document: "[Path to research document]"
     tags: [article, draft, topic-specific-tags]
     status: draft
     last_updated: [Current date in YYYY-MM-DD format]
     last_updated_by: [Author name]
     phase: writing
     ---

     # [Article Title]

     [Article content following the planned structure]

     ## Introduction

     [Engaging hook and context setting as planned]
     [Thesis introduction and reader value proposition]

     ## [Section 2 Title]

     [Content following plan guidance]
     [Integration of research findings and examples]
     [Expert quotes and supporting evidence]

     ## [Section 3 Title]

     [Continue following planned structure...]

     ## Conclusion

     [Summary and call to action as planned]
     [Future outlook and closing message]

     ---

     ## Sources and References
     [List of sources used, formatted for easy citation]
     - [Source 1 with URL and access date]
     - [Source 2 with URL and access date]

     ## Writing Metadata
     - **Plan Document**: [Path to plan document used]
     - **Research Document**: [Path to research document used]
     - **Word Count**: [Actual word count]
     - **Sections Completed**: [List of completed sections]
     ```

6. **Ensure content quality and completeness:**
   - Verify all planned sections are included and developed
   - Check that research findings are properly integrated
   - Ensure quotes and statistics are accurately represented
   - Confirm tone and style match plan specifications
   - Validate that target word count is approximately met
   - Review for logical flow and section transitions

7. **Track writing progress:**
   - Use TodoWrite to mark completed sections
   - Note any deviations from the original plan
   - Identify areas that may need additional attention in review phase
   - Document any challenges encountered during writing

8. **Sync and present writing summary:**
   - Run `humanlayer thoughts sync` if saving to thoughts directory
   - Present a summary of the completed article to the user
   - Highlight key sections and content achievements
   - Note actual word count vs. target
   - Provide preview of strongest content elements
   - Ask if they want to proceed to the review phase or need specific revisions

9. **Handle section-by-section writing:**
   - If user requests writing specific sections only:
   - Focus on requested sections while maintaining consistency
   - Ensure sections integrate well with overall article flow
   - Update TodoWrite to reflect partial completion
   - Provide guidance on remaining sections needed

## Writing Guidelines:

### Content Development
- Start each section with clear topic introduction
- Use research findings to support all major claims
- Include specific examples and case studies as planned
- Integrate expert quotes naturally into the narrative
- Maintain engaging and informative tone throughout

### Style and Flow
- Write for the specified target audience knowledge level
- Use clear, concise language appropriate to the topic
- Create smooth transitions between sections and paragraphs
- Vary sentence structure and length for readability
- Include subheadings for better content organization

### Research Integration
- Reference sources naturally without over-citation
- Use statistics and data points strategically
- Include diverse perspectives when available
- Fact-check all claims against research documents
- Maintain objectivity while presenting compelling arguments

### Engagement Techniques
- Use storytelling elements when appropriate
- Include relevant analogies or metaphors for complex concepts
- Ask rhetorical questions to engage readers
- Provide actionable insights and practical applications
- Create memorable opening and closing statements

## Important notes:
- Follow the established plan structure closely
- Focus on creating comprehensive, engaging content
- Maintain consistency in tone, style, and messaging
- Integrate research findings seamlessly into the narrative
- Write for the specified target audience throughout
- Create content that fulfills the article's stated objectives
- Don't self-edit heavily - save major revisions for review phase
- Prioritize completeness and flow over perfection in first draft
- Use research sources to support all major claims and arguments
- Create content that provides clear value to readers