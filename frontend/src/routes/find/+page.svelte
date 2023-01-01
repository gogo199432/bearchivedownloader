<script lang="ts">
    import type { ActionData } from './$types';
    export let form: ActionData;

</script>
<h1>Find entries with tags</h1>
<form method="POST">
    <label>
        Search Query
        <input style="width: 100%; height: 2rem; margin-bottom: 1rem" name="query" type="text">
    </label>
    <button style="width: 100%; height: 2rem">Search</button>
    {#if form?.incorrect}<p class="error">Query is wrongly formed. Check if tags are correct and if capitalization is correct</p>{/if}
</form>
{#if form?.success}
    {#if form?.entries.length === 0}
        <p>No entries found for this query</p>
    {:else}
    <table>
        <thead>
            <tr>
                <td>Title</td>
                <td>Labels</td>
            </tr>
        </thead>
        <tbody>
        {#each form?.entries as entry}
            <tr>
                <td><a href="/entrylist/{entry.id}">{entry.title}</a></td>
                <td>{entry.labels}</td>
            </tr>
        {/each}
        </tbody>
    </table>
    {/if}
{/if}
<h3>Example queries (a,b,c etc are any tags below): </h3>
<ul>
    <li>a</li>
    <li>(a)</li>
    <li>a AND b</li>
    <li>a AND (b OR c)</li>
</ul>
<h3>Possible operators</h3>
<ul>
    <li>AND</li>
    <li>OR</li>
    <li>XOR</li>
    <li>NOT</li>
</ul>
<h3>Possible tags</h3>
<table>
    <tbody>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Growth</b></td>
    </tr>
    <tr>
        <td>ag - Ass Expansion</td>
        <td>amaz - Amazons</td>
        <td>be - Breast Expansion</td>
    </tr>
    <tr>
        <td>belly - Belly Growth</td>
        <td>clit - Clit Expansion</td>
        <td>fat - Fattening</td>
    </tr>
    <tr>
        <td>gts - Giants (male or female)</td>
        <td>height - Height change <br>(less than shrink/gts)</td>
        <td>hg - Hair Growth</td>
    </tr>
    <tr>
        <td>immob - Immobilized</td>
        <td>inf - Inflation</td>
        <td>lab - Labia Expansion</td>
    </tr>
    <tr>
        <td>lip - Lip Expansion</td>
        <td>mg - Muscle Growth</td>
        <td>nip - Nipple Expansion</td>
    </tr>
    <tr>
        <td>pe - Penis Enlargement</td>
        <td>pr - Pregnancy</td>
        <td>tongue - Tongue TF</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Shrinking</b></td>
    </tr>
    <tr>
        <td>bred - Breast Reduction</td>
        <td>mred - Muscle Reduction</td>
        <td>pred - Penis Reduction</td>
    </tr>
    <tr>
        <td>shr - Shrinking</td>
        <td>slim - Weight Loss</td>
        <td>&nbsp;</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Body Modification</b></td>
    </tr>
    <tr>
        <td>ap - Age Progression</td>
        <td>ar - Age Regression</td>
        <td>amp - Amputation</td>
    </tr>
    <tr>
        <td>beaut - More Beautiful/Handsome</td>
        <td>bizarre - Bizarre Body Mod</td>
        <td>bwb - Boy With Boobs</td>
    </tr>
    <tr>
        <td>catg - Catgirl</td>
        <td>cloth - Change of Clothing</td>
        <td>conj - Conjoinment</td>
    </tr>
    <tr>
        <td>cowg - cowgirl</td>
        <td>hair - Hair Color/Style Change</td>
        <td>herm - Hermaphrodites/Futanari</td>
    </tr>
    <tr>
        <td>impl - Implants</td>
        <td>lac - Lactation</td>
        <td>mb - Multibreast</td>
    </tr>
    <tr>
        <td>mgen - Multiple penis/vagina</td>
        <td>ml - Multilimb</td>
        <td>nbm - Detachable body parts</td>
    </tr>
    <tr>
        <td>nogen - No genitals</td>
        <td>pierce - Piercings</td>
        <td>pusb - Pussyboys</td>
    </tr>
    <tr>
        <td>sens - Sensitivity Increase</td>
        <td>shem - Shemales/Dickgirls</td>
        <td>skin - Skin TF</td>
    </tr>
    <tr>
        <td>taur - Centaurs, Cowtaurs, etc.</td>
        <td>tat - Tattoos</td>
        <td>udder - Udder</td>
    </tr>
    <tr>
        <td>ugly - Uglier</td>
        <td>voice - Voice/Vocab change</td>
        <td>&nbsp;</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Full Transformation</b></td>
    </tr>
    <tr>
        <td>anim - Animal TF</td>
        <td>cloth2 - Turned Into Clothing</td>
        <td>costum - Magic Costume/Skinsuit</td>
    </tr>
    <tr>
        <td>f2m - Transgender (f2m)</td>
        <td>ftp - Turned into Body Part</td>
        <td>fur - Anthropomorphic</td>
    </tr>
    <tr>
        <td>fuse - Fusion/Merge</td>
        <td>human - <font color="red">anim/obj to human TF</font></td>
        <td>matl - Material TF (animate)</td>
    </tr>
    <tr>
        <td>myth - Mythological TF or Being</td>
        <td>nation - Nationality Change</td>
        <td>obj - Inanimate TF</td>
    </tr>
    <tr>
        <td>race - Race Change</td>
        <td>robot - Robot</td>
        <td>stat - Statue/Mannequin/Doll</td>
    </tr>
    <tr>
        <td>swap - Bodyswap</td>
        <td>tf - Other TF</td>
        <td>tg - Transgender (m2f)</td>
    </tr>
    <tr>
        <td>tged - Character TGed in earlier ep</td>
        <td>tgtf - TG TF in this episode</td>
        <td>twin - Twins/Clones</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Mental Transformations</b></td>
    </tr>
    <tr>
        <td>attract - Forced Attraction</td>
        <td>bimbo - Bimbo</td>
        <td>comp - Compulsion</td>
    </tr>
    <tr>
        <td>dream - Dream</td>
        <td>dumb - Less intelligent</td>
        <td>fetish - New fetish</td>
    </tr>
    <tr>
        <td>horn - Horniness</td>
        <td>ident - Identity Change</td>
        <td>illu - Illusion of a change</td>
    </tr>
    <tr>
        <td>mc - Mind Control</td>
        <td>mem - Memory Change</td>
        <td>orient - Orientation Change</td>
    </tr>
    <tr>
        <td>pers - Personality Change</td>
        <td>poss - Possession</td>
        <td>psi - Mental powers</td>
    </tr>
    <tr>
        <td>skill - Skills added or lost</td>
        <td>slave - Slavery</td>
        <td>slut - Slut/Nympho</td>
    </tr>
    <tr>
        <td>smart - More intelligent</td>
        <td>sub - Made Submissive</td>
        <td>&nbsp;</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Other Fetishes</b></td>
    </tr>
    <tr>
        <td>ab - Adult Baby</td>
        <td>bond - Bondage</td>
        <td>cross - Crossdressing</td>
    </tr>
    <tr>
        <td>exh - Exhibitionism / Nudity</td>
        <td>feed - Feeder / eating</td>
        <td>femdom - Female Domination</td>
    </tr>
    <tr>
        <td>ff - female/female (Lesbian)</td>
        <td>forg - Forced Orgasm</td>
        <td>grope - Groping</td>
    </tr>
    <tr>
        <td>mast - Masturbation</td>
        <td>messy - Messy</td>
        <td>mf - male/female</td>
    </tr>
    <tr>
        <td>mm - male/male (Gay Male)</td>
        <td>mpreg - Male Pregnancy</td>
        <td>nc - Non-Consensual</td>
    </tr>
    <tr>
        <td>norg - Can't Orgasm</td>
        <td>nurs - Nursing/Milk-Pumping</td>
        <td>poly - Polyamory / Harem / Orgy</td>
    </tr>
    <tr>
        <td>sm - Sadism / Masochism</td>
        <td>tent - Tentacles</td>
        <td>tickle - Tickling</td>
    </tr>
    <tr>
        <td>toy - Sex Toys</td>
        <td>unb - Unbirth</td>
        <td>vore - Vore</td>
    </tr>
    <tr>
        <td>voy - Voyeurism</td>
        <td>waste - Bodily Waste Products</td>
        <td>&nbsp;</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Other</b></td>
    </tr>
    <tr>
        <td>celeb - Celebrities</td>
        <td>corrupt - Corruption</td>
        <td>create - New person created</td>
    </tr>
    <tr>
        <td>drug - Drugs</td>
        <td>food - Food</td>
        <td>gap - Age Gap</td>
    </tr>
    <tr>
        <td>global - Global Change</td>
        <td>godlike - Gods / Godlike Powers</td>
        <td>humil - Humiliation</td>
    </tr>
    <tr>
        <td>inv - Invisibility</td>
        <td>job - Job/Occupation Change</td>
        <td>magic - Magic</td>
    </tr>
    <tr>
        <td>ntr - Netorare / Cheating</td>
        <td>offcam - Off-Camera change</td>
        <td>real - Reality change</td>
    </tr>
    <tr>
        <td>scene - Scene/Setting TF</td>
        <td>science - Science (mad or not)</td>
        <td>status - Social Status Change</td>
    </tr>
    <tr>
        <td>telep - Teleportation</td>
        <td>trait - Trait-Swapping</td>
        <td>ts - Time-stop</td>
    </tr>
    <tr>
        <td>undo - Return to Normal</td>
        <td>viol - Violence</td>
        <td>&nbsp;</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Mainstream Genre</b></td>
    </tr>
    <tr>
        <td>alien - Space Aliens</td>
        <td>fantasy - Fantasy World</td>
        <td>hist - Historical/Alt Hist</td>
    </tr>
    <tr>
        <td>horror - Horror</td>
        <td>humor - Humor/Comedy</td>
        <td>roman - Romance</td>
    </tr>
    <tr>
        <td>slice - Slice-of-life</td>
        <td>space - Space Travel</td>
        <td>super - Superpowers</td>
    </tr>
    <tr>
        <td>susp - Suspense</td>
        <td>time - Time Travel</td>
        <td>tragedy - Tragedy</td>
    </tr>
    <tr>
        <td bgcolor="#c0c0c0" colspan="3"><b>Episode Types</b></td>
    </tr>
    <tr>
        <td>au - Authors as Characters</td>
        <td>end - Storyline End (rare!)</td>
        <td>fan - Fanfic</td>
    </tr>
    <tr>
        <td>non_eng - Non-English Episode</td>
        <td>list - List Episode (little/no body)</td>
        <td>meta - Metafiction</td>
    </tr>
    <tr>
        <td>over - Crossover</td>
        <td>pic - Picture in episode</td>
        <td>restart - Restart Episode</td>
    </tr>
    <tr>
        <td>root - Storyline Root</td>
        <td>ubea - Main Cast (Jim/Sharon/etc.)</td>
        <td>which - "which media" list ep</td>
    </tr>
    </tbody>
</table>