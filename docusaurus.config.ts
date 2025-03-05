// Themes Config
import {themes} from "prism-react-renderer";

const theme = themes.github;
const darkTheme = themes.oneLight;

/** @type {import('@docusaurus/types').Config} */
const config = {
    title: "terraform provider otc-marketplace | Documentation",
    tagline: "Provider API reference",
    favicon: "https://usercontent.one/wp/iits-consulting.de/wp-content/uploads/2024/03/cropped-favicon-iits-2024-270x270.png",
    url: "https://iits.de",
    baseUrl: "/terraform-provider-otc-marketplace",
    onBrokenLinks: 'ignore',
    onBrokenMarkdownLinks: 'warn',
    organizationName: "iits",
    projectName: "Docs",
    plugins: [
        [
            '@docusaurus/plugin-client-redirects',
            {
                toExtensions: ['html'],
                redirects: [
                    {
                        to: '/docs/index',
                        from: ['/'],
                    },
                ],
            }
        ],
        'docusaurus-plugin-sass',
    ],
    scripts: [
        {
            src: "./src/js/wiki-version.js",
        },
    ],
    presets: [
        /** @type {import('@docusaurus/preset-classic').Options} */
        [
            '@docusaurus/preset-classic',
            ({
                docs: {
                    routeBasePath: '/docs',
                    editUrl: 'https://github.com/iits-consulting/terraform-provider-otc-marketplace/tree/main/',
                    showLastUpdateAuthor: true,
                    showLastUpdateTime: true,
                    editLocalizedFiles: true,
                    editCurrentVersion: true,
                    versions: {
                        current: {
                            label: 'v0.1.0',
                        },
                    },
                    lastVersion: 'current',
                },
                blog: {
                    showReadingTime: true,
                },
                theme: {},
            }),
        ],
    ],
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    themeConfig: {
        metadata: [{
            name: 'keywords',
            content: 'marketplace, otc, terraform, provider'
        }],
        colorMode: {
            defaultMode: 'dark',
            disableSwitch: false,
            respectPrefersColorScheme: true,
        },
        prism: {
            theme: theme,
            darkTheme: darkTheme,
            additionalLanguages: [
                "go",
                "hcl"
            ],
        },
    },
};

export default config;